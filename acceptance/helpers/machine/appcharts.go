package machine

import (
	"fmt"
	"os"

	"github.com/epinio/epinio/acceptance/helpers/proc"
	. "github.com/onsi/gomega"
)

func (m *Machine) MakeAppchartStateful(chartName string) string {
	// Create a custom chart referencing the tarball of the `standard-stateful` chart.
	// It exists in the set of releases for helm charts.
	// It is not distributed with epinio however.
	// At this point in time we use it only internally, for testing.

	tempFile := chartName + `.yaml`
	err := os.WriteFile(tempFile, []byte(fmt.Sprintf(`apiVersion: application.epinio.io/v1
kind: AppChart
metadata:
  namespace: epinio
  name: %s
  labels:
    app.kubernetes.io/component: epinio
    app.kubernetes.io/instance: default
    app.kubernetes.io/name: epinio-standard-stateful-app-chart
    app.kubernetes.io/part-of: epinio
spec:
  shortDescription: Epinio standard stateful deployment
  description: Epinio standard support chart for stateful application deployment
  helmChart: https://github.com/epinio/helm-charts/releases/download/epinio-application-stateful-0.1.21/epinio-application-stateful-0.1.21.tgz
`, chartName)), 0600)
	Expect(err).ToNot(HaveOccurred())

	out, err := proc.Kubectl("apply", "-f", tempFile)
	Expect(err).ToNot(HaveOccurred(), out)

	return tempFile
}

func (m *Machine) MakeAppchart(chartName string) string {
	tempFile := chartName + `.yaml`
	err := os.WriteFile(tempFile, []byte(fmt.Sprintf(`apiVersion: application.epinio.io/v1
kind: AppChart
metadata:
  namespace: epinio
  name: %s
  labels:
    app.kubernetes.io/component: epinio
    app.kubernetes.io/instance: default
    app.kubernetes.io/name: epinio-standard-app-chart
    app.kubernetes.io/part-of: epinio
spec:
  helmChart: fox
  settings:
    unknowntype:
      type: foofara
    badminton:
      type: integer
      minimum: hello
    maxbad:
      type: integer
      maximum: world
    fake:
      type: bool
      enum:
        - ignored
    foo:
      type: string
      minimum: ignored
    bar:
      type: string
      enum:
        - sna
        - fu
    floof:
      type: number
      minimum: '0'
    fox:
      type: integer
      maximum: '100'
    cat:
      type: number
      minimum: '0'
      maximum: '1'
`, chartName)), 0600)
	Expect(err).ToNot(HaveOccurred())

	out, err := proc.Kubectl("apply", "-f", tempFile)
	Expect(err).ToNot(HaveOccurred(), out)

	return tempFile
}

func (m *Machine) DeleteAppchart(tempFile string) {
	out, err := proc.Kubectl("delete", "-f", tempFile)
	Expect(err).ToNot(HaveOccurred(), out)
	os.Remove(tempFile)
}
