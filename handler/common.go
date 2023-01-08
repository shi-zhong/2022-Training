package handler

type PathIDModel struct {
    ID uint `uri:"ID" binding:"required"`
}

type PathStringIDModel struct {
    ID string `uri:"ID" binding:"required"`
}