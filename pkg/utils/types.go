package utils

type ctxKey int
type ctxStringKey string

const RidKey ctxKey = ctxKey(0)
const Authenticated ctxStringKey = ctxStringKey("authenticated")
