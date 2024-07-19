package dwarf

import (
	"debug/dwarf"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
)

type Content struct {
	Name string
	Addr uint64
	Size uint64
}

// SizeForDWARFVar need addr because it may in .bss section
// readMemory should only be called once
// return addr, size, type size, error
func SizeForDWARFVar(
	d *dwarf.Data,
	entry *dwarf.Entry,
	addr uint64,
	readMemory MemoryReader,
) ([]Content, uint64, error) {
	sizeOffset, ok := entry.Val(dwarf.AttrType).(dwarf.Offset)
	if !ok {
		return nil, 0, fmt.Errorf("failed to get type offset")
	}

	typ, err := d.Type(sizeOffset)
	if err != nil {
		return nil, 0, err
	}

	structTyp, ok := typ.(*dwarf.StructType)
	if ok {
		// check string
		// user can still define a struct has this name, but it's rare
		if structTyp.StructName == "string" {
			strAddr, size, err := readString(structTyp, addr, readMemory)
			if err != nil || size == 0 {
				return nil, uint64(typ.Size()), err
			}

			return []Content{{
				Name: "string",
				Addr: strAddr,
				Size: size,
			}}, uint64(typ.Size()), nil
		} else if structTyp.StructName == "[]uint8" {
			// check byte slice, normally it comes from embed
			dataAddr, size, err := readSlice(structTyp, addr, readMemory, "*uint8")
			if err != nil || size == 0 {
				return nil, uint64(typ.Size()), err
			}

			return []Content{{
				Name: "[]uint8",
				Addr: dataAddr,
				Size: size,
			}}, uint64(typ.Size()), nil
		}
	} else {
		typeDefTyp, ok := typ.(*dwarf.TypedefType)
		if ok {
			structTyp, ok = typeDefTyp.Type.(*dwarf.StructType)
			if !ok {
				return nil, uint64(typ.Size()), nil
			}

			if structTyp.StructName == "embed.FS" {
				// check embed.FS
				parts, err := readEmbedFS(structTyp, addr, readMemory)
				if err != nil || len(parts) == 0 {
					return nil, uint64(typ.Size()), err
				}

				return parts, uint64(typ.Size()), nil
			}
		}
	}

	return nil, uint64(typ.Size()), nil
}

var boolIgnores = [...]dwarf.Attr{
	dwarf.AttrCallAllCalls,
	dwarf.AttrCallAllTailCalls,
}

var ignores = [...]dwarf.Attr{
	dwarf.AttrAbstractOrigin,
	dwarf.AttrSpecification,
}

func EntryShouldIgnore(entry *dwarf.Entry) bool {
	declaration := entry.Val(dwarf.AttrDeclaration)
	if declaration != nil {
		val, ok := declaration.(bool)
		return !ok || val
	}

	inline := entry.Val(dwarf.AttrInline)
	if inline != nil {
		val, ok := inline.(int64)
		if ok && val != 0 {
			return true
		}
	}

	for _, ignore := range boolIgnores {
		valAny := entry.Val(ignore)
		if valAny != nil {
			val, ok := valAny.(bool)
			if !ok {
				slog.Warn(fmt.Sprintf("Failed to load DWARF function as bool field type unexpected %T: %s", valAny, EntryPrettyPrint(entry)))
				return true
			}
			if val {
				return true
			}
		}
	}

	for _, ignore := range ignores {
		if entry.Val(ignore) != nil {
			return true
		}
	}

	externalAny := entry.Val(dwarf.AttrExternal)
	if externalAny != nil {
		external, ok := externalAny.(bool)
		if !ok {
			slog.Debug(fmt.Sprintf("Failed to load DWARF function as dwarf.AttrExternal type unexpected %T: %s", externalAny, EntryPrettyPrint(entry)))
			return true
		}

		if external {
			if entry.Tag == dwarf.TagSubprogram {
				// external function doesn't exist in this entry
				return true
			}
		}
	}

	return false
}

func EntryFileReader(cu *dwarf.Entry, d *dwarf.Data) func(entry *dwarf.Entry) string {
	var files []*dwarf.LineFile
	lr, err := d.LineReader(cu)
	if err != nil {
		slog.Warn(fmt.Sprintf("Failed to read DWARF line: %v", err))
	}
	if lr != nil {
		files = lr.Files()
	}

	return func(entry *dwarf.Entry) string {
		const defaultName = "<autogenerated>"
		if entry.Val(dwarf.AttrTrampoline) == nil {
			fileIndexAny := entry.Val(dwarf.AttrDeclFile)
			if fileIndexAny == nil {
				slog.Debug(fmt.Sprintf("Failed to load DWARF function file as no AttrDeclFile field: %s", EntryPrettyPrint(entry)))
				return defaultName
			}
			fileIndex, ok := fileIndexAny.(int64)
			if !ok {
				slog.Warn(fmt.Sprintf("Failed to load DWARF function file as type unexpected %T: %s", fileIndexAny, EntryPrettyPrint(entry)))
				return defaultName
			}
			if fileIndex < 0 || int(fileIndex) >= len(files) {
				slog.Warn(fmt.Sprintf("Failed to load DWARF function file as index out of range %d: %s", fileIndex, EntryPrettyPrint(entry)))
				return defaultName
			}

			return files[fileIndex].Name
		}

		return defaultName
	}
}

func EntryPrettyPrint(entry *dwarf.Entry) string {
	ret := new(strings.Builder)
	ret.WriteString(entry.Tag.String())
	ret.WriteString(" ")
	ret.WriteString(strconv.Itoa(int(entry.Offset)))
	ret.WriteString(" ")
	for _, field := range entry.Field {
		ret.WriteString(fmt.Sprintf("%#v ", field))
	}

	return ret.String()
}
