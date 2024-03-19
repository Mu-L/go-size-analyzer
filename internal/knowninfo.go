package internal

import (
	"debug/elf"
	"debug/gosym"
	"debug/macho"
	"debug/pe"
	"github.com/Zxilly/go-size-analyzer/internal/tool"
	"github.com/Zxilly/go-size-analyzer/internal/wrapper"
	"github.com/goretk/gore"
	"log"
	"reflect"
	"strings"
	"unsafe"
)

type KnownInfo struct {
	Size       uint64
	BuildInfo  *gore.BuildInfo
	SectionMap *SectionMap
	Packages   *MainPackages
	KnownAddr  *KnownAddr

	gore    *gore.GoFile
	wrapper wrapper.RawFileWrapper

	VersionFlag struct {
		Leq118 bool
		Meq120 bool
	}
}

func NewKnownInfo(file *gore.GoFile) *KnownInfo {
	// ensure we have the version
	k := &KnownInfo{
		KnownAddr: NewFoundAddr(),
		Size:      tool.GetFileSize(file.GetFile()),
		BuildInfo: file.BuildInfo,

		gore:    file,
		wrapper: wrapper.NewWrapper(file.GetParsedFile()),
	}
	k.UpdateVersionFlag()
	return k
}

func (k *KnownInfo) LoadSectionMap() {
	log.Println("Loading sections...")

	sections := &SectionMap{Sections: make(map[string]*Section)}

	switch f := k.gore.GetParsedFile().(type) {
	case *pe.File:
		sections.loadFromPe(f)
	case *elf.File:
		sections.loadFromElf(f)
	case *macho.File:
		sections.loadFromMacho(f)
	default:
		panic("unreachable")
	}

	log.Println("Loading sections done")

	k.SectionMap = sections

	return
}

func (k *KnownInfo) AnalyzeSymbol(file *gore.GoFile) error {
	log.Println("Analyzing symbols...")
	var err error

	switch f := file.GetParsedFile().(type) {
	case *pe.File:
		err = analyzePeSymbol(f, k)
	case *elf.File:
		err = analyzeElfSymbol(f, k)
	case *macho.File:
		err = analyzeMachoSymbol(f, k)
	default:
		panic("unreachable")
	}

	if err != nil {
		return err
	}

	k.KnownAddr.BuildSymbolCoverage()

	log.Println("Analyzing symbols done")

	return nil
}

func (k *KnownInfo) Validate() error {
	// TODO: validate the result
	return nil
}

func (k *KnownInfo) UpdateVersionFlag() {
	ver, err := k.gore.GetCompilerVersion()
	if err != nil {
		// if we can't get build info, we assume it's go1.20 plus
		k.VersionFlag.Meq120 = true
	} else {
		k.VersionFlag.Leq118 = gore.GoVersionCompare(ver.Name, "go1.18.10") <= 0
		k.VersionFlag.Meq120 = gore.GoVersionCompare(ver.Name, "go1.20rc1") >= 0
	}
}

// ExtractPackageFromSymbol copied from debug/gosym/symtab.go
func (k *KnownInfo) ExtractPackageFromSymbol(s string) string {
	sym := &gosym.Sym{
		Name: s,
	}

	val := reflect.ValueOf(sym).Elem()
	ver := val.FieldByName("goVersion")

	set := func(i int) {
		reflect.NewAt(ver.Type(), unsafe.Pointer(ver.UnsafeAddr())).Elem().SetInt(int64(i))
	}

	if k.VersionFlag.Meq120 {
		set(5) // ver120
	} else if k.VersionFlag.Leq118 {
		set(4) // ver118
	}

	pn := sym.PackageName()

	if strings.Count(pn, ".") >= 3 {
		// see MainPackages.Add
		return ""
	}
	return pn
}

func (k *KnownInfo) GetPaddingSize() uint64 {
	var sectionSize uint64 = 0
	for _, section := range k.SectionMap.Sections {
		sectionSize += section.Size
	}
	return k.Size - sectionSize
}

func (k *KnownInfo) RequireModInfo() {
	if k.BuildInfo == nil || len(k.BuildInfo.ModInfo.Deps) == 0 {
		log.Fatal("mod info is required for this operation")
	}
}