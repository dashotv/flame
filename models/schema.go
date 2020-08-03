package models

type GolemSchema struct {
    Download *DownloadSchema
    
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
    
}
