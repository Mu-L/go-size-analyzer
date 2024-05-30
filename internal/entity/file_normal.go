//go:build !js && !wasm

package entity

import (
	"github.com/Zxilly/go-size-analyzer/internal/utils"
	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

var FileMarshalerCompact = json.MarshalFuncV2[File](func(encoder *jsontext.Encoder, file File, options json.Options) error {
	utils.Must(encoder.WriteToken(jsontext.ObjectStart))

	utils.Must(json.MarshalEncode(encoder, "file_path", options))
	utils.Must(json.MarshalEncode(encoder, file.FilePath, options))
	utils.Must(json.MarshalEncode(encoder, "size", options))
	utils.Must(json.MarshalEncode(encoder, file.FullSize(), options))
	utils.Must(json.MarshalEncode(encoder, "pcln_size", options))
	utils.Must(json.MarshalEncode(encoder, file.PclnSize(), options))

	utils.Must(encoder.WriteToken(jsontext.ObjectEnd))
	return nil
})
