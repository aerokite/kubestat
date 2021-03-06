package pkg

import (
	apps "k8s.io/api/apps/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

// getStatefulset ...
func (c *Client) getStatefulset(ns string) (*apps.StatefulSetList, error) {
	if ns == "" {
		ns = meta.NamespaceAll
	}
	return c.kubernetesclient.AppsV1().StatefulSets(ns).List(meta.ListOptions{
		LabelSelector: labels.Everything().String(),
	})
}

// GetStatefulsetMetrics ...
func (c *Client) GetStatefulsetMetrics(ns string) (NamespaceWiseServiceMetrics, error) {
	statefulsets, err := c.getStatefulset(ns)
	if err != nil {
		return nil, err
	}

	metrics := make(NamespaceWiseServiceMetrics)

	for _, deploy := range statefulsets.Items {
		selector, err := meta.LabelSelectorAsSelector(deploy.Spec.Selector)
		if err != nil {
			return nil, err
		}
		podMetricsList, err := c.metricsclient.Metrics().PodMetricses(deploy.Namespace).List(meta.ListOptions{
			LabelSelector: selector.String(),
		})
		if err != nil {
			return nil, err
		}

		metric := processServiceWiseMetrics(podMetricsList)
		if _, found := metrics[deploy.Namespace]; !found {
			metrics[deploy.Namespace] = make(map[string]*Metric)
		}

		metrics[deploy.Namespace][deploy.Name] = metric

	}

	return metrics, nil
}
