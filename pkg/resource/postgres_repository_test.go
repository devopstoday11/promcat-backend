package resource_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/falcosecurity/cloud-native-security-hub/test/fixtures/resources"

	"database/sql"
	"os"

	"github.com/falcosecurity/cloud-native-security-hub/pkg/resource"
)

var _ = Describe("Postgres Resource Repository", func() {
	var repository resource.Repository

	BeforeEach(func() {
		db, _ := sql.Open("postgres", os.Getenv("DATABASE_URL"))
		repository = resource.NewPostgresRepository(db)

		db.Exec("TRUNCATE TABLE security_resources")
		db.Exec("TRUNCATE TABLE latest_security_resources")
	})

	It("saves a new resource", func() {
		repository.Save(resources.Apache())

		retrieved, _ := repository.FindById("apache")
		Expect(retrieved).To(Equal(resources.Apache()))
	})

	It("retrieves all existent resources", func() {
		repository.Save(resources.Apache())
		repository.Save(resources.MongoDB())

		retrieved, _ := repository.FindAll()

		Expect(retrieved).To(Equal([]*resource.Resource{
			resources.Apache(),
			resources.MongoDB()}))
	})

	Context("when querying by id", func() {
		Context("and resource is not found", func() {
			It("returns an error", func() {
				retrieved, err := repository.FindById("non existent id")

				Expect(retrieved).To(BeNil())
				Expect(err).To(HaveOccurred())
			})
		})

		It("returns latest version of the resource", func() {
			apache := resources.Apache()
			repository.Save(apache)

			apache.Version = "2.0.0"
			repository.Save(apache)

			retrieved, _ := repository.FindById("apache")

			Expect(retrieved).To(Equal(apache))
		})

		Context("and version is specified as well", func() {
			It("returns the resource with the specified version", func() {
				apache := resources.Apache()
				repository.Save(apache)

				apache.Version = "2.0.0"
				repository.Save(apache)

				retrieved, _ := repository.FindByVersion("apache", "1.0.0")

				Expect(retrieved).To(Equal(resources.Apache()))
			})
		})
	})

	Context("when saving several versions for a resource", func() {
		It("returns all all available versions", func() {
			apache := resources.Apache()
			repository.Save(apache)

			apache.Version = "2.0.0"
			repository.Save(apache)

			retrieved, _ := repository.FindById("apache")

			Expect(retrieved.AvailableVersions).To(Equal([]string{"1.0.0", "2.0.0"}))
		})
	})
})
