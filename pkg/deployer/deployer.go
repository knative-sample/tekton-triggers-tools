package manager

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"knative.dev/serving/pkg/apis/config"
	"knative.dev/serving/pkg/apis/serving/v1"
	servingclient "knative.dev/serving/pkg/client/injection/client"
	kserviceinformer "knative.dev/serving/pkg/client/injection/informers/serving/v1/service"

	"k8s.io/apimachinery/pkg/api/errors"
	"knative.dev/pkg/logging"
)

func NewDeployer(
	ctx context.Context,
) *Deployer {
	logger := logging.FromContext(ctx)
	serviceInformer := kserviceinformer.Get(ctx)

	m := &Deployer{
		logger:          logger,
		servingclient:   servingclient.Get(ctx),
		serviceInformer: serviceInformer,
	}

	return m
}

func (m *Deployer) Deploy(ksvcname, namespace string, image string) error {
	logger := m.logger
	logger.Infof("Deployer is run")

	ksvc, err := m.servingclient.ServingV1().Services(namespace).Get(ksvcname, metav1.GetOptions{})
	if err != nil {
		if !errors.IsNotFound(err) {
			logger.Fatalf("get ksvc %s/%s info error:%s", namespace, ksvcname, err.Error())
		}

		// 这是一个创建 ksvc 的例子
		if err := m.CreateKsvc(ksvcname, namespace, image); err != nil {
			logger.Fatalf("create ksvc error:%s", err.Error())
			return err
		}
		logger.Infof("create ksvc %s/%s image:%s success", namespace, ksvcname, image)
		return nil
	}

	for index, c := range ksvc.Spec.Template.Spec.Containers {
		if c.Name == config.DefaultUserContainerName {
			ksvc.Spec.Template.Spec.Containers[index].Image = image
		}
	}

	if _, err := m.servingclient.ServingV1().Services(namespace).Update(ksvc); err != nil {
		// TODO 对 etcd Revision 有更新的情况进行兼容
		logger.Fatalf("get ksvc %s/%s image:%s error:%s", namespace, ksvcname, image, err.Error())
	}
	logger.Infof("update ksvc %s/%s image:%s success", namespace, ksvcname, image)

	return nil
}

func (m *Deployer) CreateKsvc(ksvcname, namespace string, image string) error {
	logger := m.logger
	ksvc := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ksvcname,
			Namespace: namespace,
		},
		Spec: v1.ServiceSpec{
			ConfigurationSpec: v1.ConfigurationSpec{
				Template: v1.RevisionTemplateSpec{
					Spec: v1.RevisionSpec{
						PodSpec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Image: image,
								},
							},
						},
					},
				},
			},
		},
	}
	_, err := m.servingclient.ServingV1().Services("default").Create(ksvc)
	if err != nil {
		logger.Errorf("create ksvc %s/%s image:%s error:%s", namespace, ksvcname, image, err.Error())
		return err
	}

	return nil
}
