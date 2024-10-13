package validateadmissionwebhook

import (
	"net/http"
	"statefulset/base"
	"statefulset/cmds/server/context"
	"strings"

	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/internalversion/scheme"
)

var (
	whitelistedRegistries = []string{"docker.io", "gcr.io"}
)

func ValidateAdmission(c *context.Context) error {
	base.Context.Logger.Debugf("handler ValidateAdmission")

	deserializer := scheme.Codecs.UniversalDeserializer()
	var res admissionv1.AdmissionResponse
	var req admissionv1.AdmissionReview
	if err := c.ParseJSONBody(req, false); err != nil {
		base.Context.Logger.Errorf("[ValidateAdmission] parse body failed, err:%v", err)
		res.Allowed = false
		res.Result.Message = "unexpected type"
		res.Result.Reason = "pares body failed"
	}

	pod := &corev1.Pod{}
	if _, _, err := deserializer.Decode(req.Request.Object.Raw, nil, pod); err != nil {
		res.Allowed = false
		res.Result.Message = "unexpected type"
		res.Result.Reason = "BadRequest"
	}

	for _, container := range pod.Spec.Containers {
		for _, registry := range whitelistedRegistries {
			if strings.HasPrefix(container.Image, registry) {
				continue
			}
			res.Allowed = false
			res.Result.Message = "illegal registries"
			res.Result.Reason = "BadRequest"
		}
	}
	c.WriteJSONResponse(http.StatusOK, res)
	return nil
}
