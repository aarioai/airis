package entity

import (
	"github.com/aarioai/airis-driver/driver/index"
	"github.com/aarioai/airis/aa/atype"
)

type SqlDemo struct {
	// SQL demo:
	Id            uint64         `db:"id"`
	Content       atype.Text     `db:"text"`
	Image         atype.Image    `db:"image"`
	Images        atype.Images   `db:"images"`
	Audio         atype.Audio    `db:"audio"`
	Audios        atype.Audios   `db:"audio"`
	Video         atype.Video    `db:"video"`
	Videos        atype.Videos   `db:"video"`
	Position      atype.Position `db:"pos" options:"no_update"`
	EffectiveDate atype.Date     `db:"effective_date"  options:"no_update"`
	CreatedAt     atype.Datetime `db:"created_at"`
	UpdatedAt     atype.Datetime `db:"updated_at"`
}

func (t SqlDemo) Table() string {
	return "law"
}

func (t SqlDemo) Indexes() index.Indexes {
	return index.Indexes{
		"PRIMARY":   index.Primary("id"),
		"k_image":   index.Index("image"),
		"s_pos":     index.Spatial2D("pos"),
		"f_content": index.FullText("content"),
	}
}

// same as above
func (t SqlDemo) Indexes() index.Indexes {
	return index.NewIndexes(
		index.Primary("id"),
		index.Index("image"),
		index.Spatial2D("pos"),
		index.FullText("content"),
	)
}
