package fastreq

type Middleware func(ctx *Ctx) error
