package models

type GolemSchema struct {
    Download *DownloadSchema
    Medium *MediumSchema
    Release *ReleaseSchema
    
}

type DownloadSchema struct {
    MediumId *FieldSchema
    Auto *FieldSchema
    Multi *FieldSchema
    Force *FieldSchema
    Url *FieldSchema
    ReleaseId *FieldSchema
    Thash *FieldSchema
    Timestamps *FieldSchema
    Selected *FieldSchema
    Status *FieldSchema
    Files *FieldSchema
    
}
type MediumSchema struct {
    Type *FieldSchema
    Kind *FieldSchema
    Source *FieldSchema
    SourceId *FieldSchema
    Title *FieldSchema
    Description *FieldSchema
    Slug *FieldSchema
    Text *FieldSchema
    Display *FieldSchema
    Directory *FieldSchema
    Search *FieldSchema
    SearchParams *FieldSchema
    Active *FieldSchema
    Downloaded *FieldSchema
    Completed *FieldSchema
    Skipped *FieldSchema
    Watched *FieldSchema
    Broken *FieldSchema
    ReleaseDate *FieldSchema
    Paths *FieldSchema
    
}
type ReleaseSchema struct {
    Type *FieldSchema
    Source *FieldSchema
    Raw *FieldSchema
    Title *FieldSchema
    Description *FieldSchema
    Size *FieldSchema
    View *FieldSchema
    Download *FieldSchema
    Infohash *FieldSchema
    Name *FieldSchema
    Season *FieldSchema
    Episode *FieldSchema
    Volume *FieldSchema
    Checksum *FieldSchema
    Group *FieldSchema
    Author *FieldSchema
    Verified *FieldSchema
    Widescreen *FieldSchema
    Uncensored *FieldSchema
    Bluray *FieldSchema
    Resolution *FieldSchema
    Encoding *FieldSchema
    Quality *FieldSchema
    Published *FieldSchema
    
}

type FieldSchema struct {
    Name string
    Type string
}

var Schema = &GolemSchema {
    Download: &DownloadSchema {
        MediumId: &FieldSchema{
            Name: "medium_id",
            Type: "primitive.ObjectID",
        },
        Auto: &FieldSchema{
            Name: "auto",
            Type: "bool",
        },
        Multi: &FieldSchema{
            Name: "multi",
            Type: "bool",
        },
        Force: &FieldSchema{
            Name: "force",
            Type: "bool",
        },
        Url: &FieldSchema{
            Name: "url",
            Type: "string",
        },
        ReleaseId: &FieldSchema{
            Name: "release_id",
            Type: "string",
        },
        Thash: &FieldSchema{
            Name: "thash",
            Type: "string",
        },
        Timestamps: &FieldSchema{
            Name: "timestamps",
            Type: "struct",
        },
        Selected: &FieldSchema{
            Name: "selected",
            Type: "string",
        },
        Status: &FieldSchema{
            Name: "status",
            Type: "string",
        },
        Files: &FieldSchema{
            Name: "files",
            Type: "[]struct",
        },
        
    },
    Medium: &MediumSchema {
        Type: &FieldSchema{
            Name: "type",
            Type: "string",
        },
        Kind: &FieldSchema{
            Name: "kind",
            Type: "primitive.Symbol",
        },
        Source: &FieldSchema{
            Name: "source",
            Type: "string",
        },
        SourceId: &FieldSchema{
            Name: "source_id",
            Type: "string",
        },
        Title: &FieldSchema{
            Name: "title",
            Type: "string",
        },
        Description: &FieldSchema{
            Name: "description",
            Type: "string",
        },
        Slug: &FieldSchema{
            Name: "slug",
            Type: "string",
        },
        Text: &FieldSchema{
            Name: "text",
            Type: "[]string",
        },
        Display: &FieldSchema{
            Name: "display",
            Type: "string",
        },
        Directory: &FieldSchema{
            Name: "directory",
            Type: "string",
        },
        Search: &FieldSchema{
            Name: "search",
            Type: "string",
        },
        SearchParams: &FieldSchema{
            Name: "search_params",
            Type: "struct",
        },
        Active: &FieldSchema{
            Name: "active",
            Type: "bool",
        },
        Downloaded: &FieldSchema{
            Name: "downloaded",
            Type: "bool",
        },
        Completed: &FieldSchema{
            Name: "completed",
            Type: "bool",
        },
        Skipped: &FieldSchema{
            Name: "skipped",
            Type: "bool",
        },
        Watched: &FieldSchema{
            Name: "watched",
            Type: "bool",
        },
        Broken: &FieldSchema{
            Name: "broken",
            Type: "bool",
        },
        ReleaseDate: &FieldSchema{
            Name: "release_date",
            Type: "time.Time",
        },
        Paths: &FieldSchema{
            Name: "paths",
            Type: "[]struct",
        },
        
    },
    Release: &ReleaseSchema {
        Type: &FieldSchema{
            Name: "type",
            Type: "string",
        },
        Source: &FieldSchema{
            Name: "source",
            Type: "string",
        },
        Raw: &FieldSchema{
            Name: "raw",
            Type: "string",
        },
        Title: &FieldSchema{
            Name: "title",
            Type: "string",
        },
        Description: &FieldSchema{
            Name: "description",
            Type: "string",
        },
        Size: &FieldSchema{
            Name: "size",
            Type: "string",
        },
        View: &FieldSchema{
            Name: "view",
            Type: "string",
        },
        Download: &FieldSchema{
            Name: "download",
            Type: "string",
        },
        Infohash: &FieldSchema{
            Name: "infohash",
            Type: "string",
        },
        Name: &FieldSchema{
            Name: "name",
            Type: "string",
        },
        Season: &FieldSchema{
            Name: "season",
            Type: "int",
        },
        Episode: &FieldSchema{
            Name: "episode",
            Type: "int",
        },
        Volume: &FieldSchema{
            Name: "volume",
            Type: "int",
        },
        Checksum: &FieldSchema{
            Name: "checksum",
            Type: "string",
        },
        Group: &FieldSchema{
            Name: "group",
            Type: "string",
        },
        Author: &FieldSchema{
            Name: "author",
            Type: "string",
        },
        Verified: &FieldSchema{
            Name: "verified",
            Type: "bool",
        },
        Widescreen: &FieldSchema{
            Name: "widescreen",
            Type: "bool",
        },
        Uncensored: &FieldSchema{
            Name: "uncensored",
            Type: "bool",
        },
        Bluray: &FieldSchema{
            Name: "bluray",
            Type: "bool",
        },
        Resolution: &FieldSchema{
            Name: "resolution",
            Type: "string",
        },
        Encoding: &FieldSchema{
            Name: "encoding",
            Type: "string",
        },
        Quality: &FieldSchema{
            Name: "quality",
            Type: "string",
        },
        Published: &FieldSchema{
            Name: "published",
            Type: "time.Time",
        },
        
    },
    
}
