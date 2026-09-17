package assemble

import (
	"github.com/MyCode83/godirb/internal/core"
	"github.com/MyCode83/godirb/internal/signature"
)

func AttachSignatures(engine *core.Core) error {
	signatures, err := signature.New()
	if err != nil {
		return err
	}

	engine.Signatures = signatures
	return nil
}
