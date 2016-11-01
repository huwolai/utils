package openshift

const (
	DEFUALT_NAMESPACE = "qiyunxin"
	PROTOCOL  =  "http"
	FIX_DOMAIN  = "svc.cluster.local"
	//服务
	SERVICE_PORT = "8080"
	//应用管理
	APPMANAGER_PORT = "8081"
	//权限管理
	SECURITYMANAGER_PORT = "8082"
)
//获取权限管理服务地址
func GetServiceSecurityUrl(serviceId string) string  {

	return GetServiceSecurityUrlWithNamespace(DEFUALT_NAMESPACE,serviceId)

}

//获取权限管理服务地址
func GetServiceSecurityUrlWithNamespace(namespace,serviceId string) string  {

	return GetServiceUrl(namespace,serviceId,SECURITYMANAGER_PORT)
}

//获取APP服务地址
func GetServiceAppUrl(serviceId string)  string {

	return GetServiceAppUrlWithNamespace(DEFUALT_NAMESPACE,serviceId)
}

//通过空间获取APP服务地址
func GetServiceAppUrlWithNamespace(namespace,serviceId string) string  {

	return GetServiceUrl(namespace,serviceId,APPMANAGER_PORT)
}
//获取服务的URL
func GetServiceUrl(namespace,serviceId ,port string) string {

	return PROTOCOL +"://" + serviceId+"." +namespace +"." + FIX_DOMAIN + ":" + port
}
