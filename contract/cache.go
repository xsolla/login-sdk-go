package contract

type ValidationKeysCache interface {
	Get(projectId string) (interface{}, bool)
	Set(projectId string, keys interface{})
}
