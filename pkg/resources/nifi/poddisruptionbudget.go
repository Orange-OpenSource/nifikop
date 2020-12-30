package nifi

import (
	"fmt"
	"math"
	runtimeClient "sigs.k8s.io/controller-runtime/pkg/client"
	"strconv"
	"strings"

	"github.com/Orange-OpenSource/nifikop/pkg/resources/templates"
	"github.com/Orange-OpenSource/nifikop/pkg/util"
	"github.com/go-logr/logr"
	policyv1beta1 "k8s.io/api/policy/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func (r *Reconciler) podDisruptionBudget(log logr.Logger) (runtimeClient.Object, error) {
	minAvailable, err := r.computeMinAvailable(log)

	if err != nil {
		return nil, err

	}

	return &policyv1beta1.PodDisruptionBudget{
		TypeMeta: metav1.TypeMeta{
			Kind:       "PodDisruptionBudget",
			APIVersion: "policy/v1beta1",
		},
		ObjectMeta: templates.ObjectMetaWithAnnotations(
			fmt.Sprintf("%s-pdb", r.NifiCluster.Name),
			util.MergeLabels(LabelsForNifi(r.NifiCluster.Name), r.NifiCluster.Labels),
			r.NifiCluster.Spec.Service.Annotations,
			r.NifiCluster,
		),
		Spec: policyv1beta1.PodDisruptionBudgetSpec{
			MinAvailable: &minAvailable,
			Selector: &metav1.LabelSelector{
				MatchLabels: LabelsForNifi(r.NifiCluster.Name),
			},
		},
	}, nil

}

// Calculate maxUnavailable as max between nodeCount - 1 (so we only allow 1 node to be disrupted)
// and 1 (to cover for 1 node clusters)
func (r *Reconciler) computeMinAvailable(log logr.Logger) (intstr.IntOrString, error) {

	/*
		budget = r.KafkaCluster.Spec.DisruptionBudget.budget (string) ->
		- can either be %percentage or static number
		Logic:
		Max(1, brokers-budget) - for a static number budget
		Max(1, brokers-brokers*percentage) - for a percentage budget
	*/

	// number of brokers in the NifiCluster
	nodes := len(r.NifiCluster.Spec.Nodes)

	// configured budget in the NifiCluster
	disruptionBudget := r.NifiCluster.Spec.DisruptionBudget.Budget

	budget := 0

	// treat percentage budget
	if strings.HasSuffix(disruptionBudget, "%") {
		percentage, err := strconv.ParseFloat(disruptionBudget[:len(disruptionBudget)-1], 4)
		if err != nil {
			log.Error(err, "error occurred during parsing the disruption budget")
			return intstr.FromInt(-1), err
		} else {
			budget = int(math.Floor((percentage * float64(nodes)) / 100))
		}
	} else {
		// treat static number budget
		staticBudget, err := strconv.ParseInt(disruptionBudget, 10, 0)
		if err != nil {
			log.Error(err, "error occurred during parsing the disruption budget")
			return intstr.FromInt(-1), err
		} else {
			budget = int(staticBudget)
		}

	}

	return intstr.FromInt(util.Max(1, nodes-budget)), nil
}
