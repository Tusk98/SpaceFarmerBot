package booru

type BooruPost struct {
    ID int
    ImageWidth int
    ImageHeight int
	PreviewFileUrl string
    FileUrl string
    TagsGeneral string
    TagsCharacter string
    TagsCopyright string
    TagsArtist string
    TagsMeta string
}
