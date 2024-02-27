package dml

const captureInsertKeyspaceAndTable = `(?is)INSERT\s+INTO\s+([A-Za-z_0-9.]+)`
const captureInsertReturnValues = `(?is)RETURNING\s+(.*)`

// TODO
const captureInsertSingleInputValues = `(?i)([A-Za-z_0-9()]+)\s*=\s*\?`
const captureInsertSliceInputValues = `(?i)([A-Za-z_0-9()]+)\s+IN\s*\(%s`
