package controller

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	instancev1alpha1 "github.com/prot0s34/neovim-operator/api/v1alpha1"
)

// NeovimReconciler reconciles a Neovim object
type NeovimReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=instance.neovim.prot0s.com,resources=neovims,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=instance.neovim.prot0s.com,resources=neovims/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=instance.neovim.prot0s.com,resources=neovims/finalizers,verbs=update
//+kubebuilder:rbac:groups="",resources=pods,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=pods/log,verbs=get;list
//+kubebuilder:rbac:groups="",resources=pods/exec,verbs=create;get;list

func (r *NeovimReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	neovim := &instancev1alpha1.Neovim{}
	err := r.Get(ctx, req.NamespacedName, neovim)
	if err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		log.Error(err, "Failed to get Neovim Instance Object")
		return ctrl.Result{}, err
	}

	pod := &corev1.Pod{}
	err = r.Get(ctx, types.NamespacedName{Name: neovim.Name, Namespace: neovim.Namespace}, pod)
	if err != nil && errors.IsNotFound(err) {
		newPod := constructPodForNeovim(neovim)
		log.Info("Creating a new Pod", "Pod.Namespace", newPod.Namespace, "Pod.Name", newPod.Name)
		err = r.Create(ctx, newPod)
		if err != nil {
			log.Error(err, "Failed to create new Pod", "Pod.Namespace", newPod.Namespace, "Pod.Name", newPod.Name)
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		log.Error(err, "Failed to get Pod")
		return ctrl.Result{}, err
	}

	// Pod already exists - no need to requeue
	log.Info("Skip reconcile: Pod already exists", "Pod.Namespace", pod.Namespace, "Pod.Name", pod.Name)
	return ctrl.Result{}, nil
}

// constructPodForNeovim creates a pod for a Neovim resource
func constructPodForNeovim(neovim *instancev1alpha1.Neovim) *corev1.Pod {
	labels := map[string]string{
		"app": neovim.Name,
	}

	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      neovim.Name,
			Namespace: neovim.Namespace,
			Labels:    labels,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "neovim",
					Image: "mashmb/nvim:dev",
				},
			},
		},
	}
}

func (r *NeovimReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&instancev1alpha1.Neovim{}).
		Complete(r)
}
