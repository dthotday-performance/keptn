package handler

import (
	"encoding/json"
	"errors"
	keptnapimodels "github.com/keptn/go-utils/pkg/api/models"
	db_mock "github.com/keptn/keptn/shipyard-controller/db/mock"
	"github.com/keptn/keptn/shipyard-controller/models"
	"github.com/keptn/keptn/shipyard-controller/operations"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetProjects(t *testing.T) {
	secretStore := &SecretStoreMock{}
	projectRepo := &db_mock.ProjectRepoMock{}
	eventRepo := &db_mock.EventRepoMock{}
	taskSequenceRepo := &db_mock.TaskSequenceRepoMock{}
	configStore := &ConfigurationStoreMock{}

	p1 := &models.ExpandedProject{}
	p2 := &models.ExpandedProject{}
	expectedProjects := []*models.ExpandedProject{p1, p2}

	projectRepo.GetProjectsFunc = func() ([]*models.ExpandedProject, error) {
		return expectedProjects, nil
	}

	instance := NewProjectManager(configStore, secretStore, projectRepo, taskSequenceRepo, eventRepo)
	actualProjects, err := instance.Get()
	assert.Nil(t, err)
	assert.Equal(t, expectedProjects, actualProjects)
}

func TestGetProjectsErr(t *testing.T) {
	secretStore := &SecretStoreMock{}
	projectRepo := &db_mock.ProjectRepoMock{}
	eventRepo := &db_mock.EventRepoMock{}
	taskSequenceRepo := &db_mock.TaskSequenceRepoMock{}
	configStore := &ConfigurationStoreMock{}

	projectRepo.GetProjectsFunc = func() ([]*models.ExpandedProject, error) {
		return nil, errors.New("Oh Oh...")
	}

	instance := NewProjectManager(configStore, secretStore, projectRepo, taskSequenceRepo, eventRepo)
	actualProjects, err := instance.Get()
	assert.NotNil(t, err)
	assert.Nil(t, actualProjects)
}

func TestGetByName(t *testing.T) {
	secretStore := &SecretStoreMock{}
	projectRepo := &db_mock.ProjectRepoMock{}
	eventRepo := &db_mock.EventRepoMock{}
	taskSequenceRepo := &db_mock.TaskSequenceRepoMock{}
	configStore := &ConfigurationStoreMock{}

	projectRepo.GetProjectFunc = func(projectName string) (*models.ExpandedProject, error) {
		return &models.ExpandedProject{}, nil
	}

	instance := NewProjectManager(configStore, secretStore, projectRepo, taskSequenceRepo, eventRepo)
	project, err := instance.GetByName("my-project")
	assert.Nil(t, err)
	assert.NotNil(t, project)
	assert.Equal(t, "my-project", projectRepo.GetProjectCalls()[0].ProjectName)
}

func TestGetByNameErr(t *testing.T) {
	secretStore := &SecretStoreMock{}
	projectRepo := &db_mock.ProjectRepoMock{}
	eventRepo := &db_mock.EventRepoMock{}
	taskSequenceRepo := &db_mock.TaskSequenceRepoMock{}
	configStore := &ConfigurationStoreMock{}

	projectRepo.GetProjectFunc = func(projectName string) (*models.ExpandedProject, error) {
		return nil, errors.New("Oh Oh...")
	}

	instance := NewProjectManager(configStore, secretStore, projectRepo, taskSequenceRepo, eventRepo)
	project, err := instance.GetByName("my-project")
	assert.NotNil(t, err)
	assert.Nil(t, project)
	assert.Equal(t, "my-project", projectRepo.GetProjectCalls()[0].ProjectName)
}

func TestCreate_GettingProjectFails(t *testing.T) {
	secretStore := &SecretStoreMock{}
	projectRepo := &db_mock.ProjectRepoMock{}
	eventRepo := &db_mock.EventRepoMock{}
	taskSequenceRepo := &db_mock.TaskSequenceRepoMock{}
	configStore := &ConfigurationStoreMock{}

	configStore.GetProjectFunc = func(in1 string) (*keptnapimodels.Project, error) {
		return nil, errors.New("Whoops...")
	}

	instance := NewProjectManager(configStore, secretStore, projectRepo, taskSequenceRepo, eventRepo)
	params := &operations.CreateProjectParams{
		GitRemoteURL: "git-url",
		GitToken:     "git-token",
		GitUser:      "git-user",
		Name:         stringp("existing-project"),
		Shipyard:     stringp("shipyard"),
	}
	err, rollback := instance.Create(params)
	assert.NotNil(t, err)
	rollback()

}

func TestCreateWithAlreadyExistingProject(t *testing.T) {
	secretStore := &SecretStoreMock{}
	projectRepo := &db_mock.ProjectRepoMock{}
	eventRepo := &db_mock.EventRepoMock{}
	taskSequenceRepo := &db_mock.TaskSequenceRepoMock{}
	configStore := &ConfigurationStoreMock{}

	configStore.GetProjectFunc = func(in1 string) (*keptnapimodels.Project, error) {
		project := &keptnapimodels.Project{
			ProjectName: "existing-project",
		}
		return project, nil
	}

	instance := NewProjectManager(configStore, secretStore, projectRepo, taskSequenceRepo, eventRepo)
	params := &operations.CreateProjectParams{
		GitRemoteURL: "git-url",
		GitToken:     "git-token",
		GitUser:      "git-user",
		Name:         stringp("existing-project"),
		Shipyard:     stringp("shipyard"),
	}
	err, rollback := instance.Create(params)
	assert.NotNil(t, err)
	rollback()

}

func TestCreate_WhenCreatingProjectInConfigStoreFails_ThenSecretGetsDeletedAgain(t *testing.T) {
	secretStore := &SecretStoreMock{}
	projectRepo := &db_mock.ProjectRepoMock{}
	eventRepo := &db_mock.EventRepoMock{}
	taskSequenceRepo := &db_mock.TaskSequenceRepoMock{}
	configStore := &ConfigurationStoreMock{}

	configStore.GetProjectFunc = func(in1 string) (*keptnapimodels.Project, error) {
		return nil, nil
	}
	configStore.CreateProjectFunc = func(keptnapimodels.Project) error {
		return errors.New("whoops...")
	}

	secretStore.UpdateSecretFunc = func(name string, content map[string][]byte) error {
		return nil
	}

	secretStore.DeleteSecretFunc = func(name string) error {
		return nil
	}

	instance := NewProjectManager(configStore, secretStore, projectRepo, taskSequenceRepo, eventRepo)
	params := &operations.CreateProjectParams{
		Name: stringp("my-project"),
	}
	err, rollback := instance.Create(params)
	assert.NotNil(t, err)
	rollback()
	assert.Equal(t, "git-credentials-my-project", secretStore.DeleteSecretCalls()[0].Name)
}

func TestCreate_WhenUploadingShipyardFails_thenProjectAndSecretGetDeletedAgain(t *testing.T) {

	encodedShipyard := "YXBpVmVyc2lvbjogc3BlYy5rZXB0bi5zaC8wLjIuMApraW5kOiBTaGlweWFyZAptZXRhZGF0YToKICBuYW1lOiB0ZXN0LXNoaXB5YXJkCnNwZWM6CiAgc3RhZ2VzOgogIC0gbmFtZTogZGV2CiAgICBzZXF1ZW5jZXM6CiAgICAtIG5hbWU6IGFydGlmYWN0LWRlbGl2ZXJ5CiAgICAgIHRhc2tzOgogICAgICAtIG5hbWU6IGRlcGxveW1lbnQKICAgICAgICBwcm9wZXJ0aWVzOiAgCiAgICAgICAgICBzdHJhdGVneTogZGlyZWN0CiAgICAgIC0gbmFtZTogdGVzdAogICAgICAgIHByb3BlcnRpZXM6CiAgICAgICAgICBraW5kOiBmdW5jdGlvbmFsCiAgICAgIC0gbmFtZTogZXZhbHVhdGlvbiAKICAgICAgLSBuYW1lOiByZWxlYXNlIAoKICAtIG5hbWU6IGhhcmRlbmluZwogICAgc2VxdWVuY2VzOgogICAgLSBuYW1lOiBhcnRpZmFjdC1kZWxpdmVyeQogICAgICB0cmlnZ2VyczoKICAgICAgLSBkZXYuYXJ0aWZhY3QtZGVsaXZlcnkuZmluaXNoZWQKICAgICAgdGFza3M6CiAgICAgIC0gbmFtZTogZGVwbG95bWVudAogICAgICAgIHByb3BlcnRpZXM6IAogICAgICAgICAgc3RyYXRlZ3k6IGJsdWVfZ3JlZW5fc2VydmljZQogICAgICAtIG5hbWU6IHRlc3QKICAgICAgICBwcm9wZXJ0aWVzOiAgCiAgICAgICAgICBraW5kOiBwZXJmb3JtYW5jZQogICAgICAtIG5hbWU6IGV2YWx1YXRpb24KICAgICAgLSBuYW1lOiByZWxlYXNlCiAgICAgICAgCiAgLSBuYW1lOiBwcm9kdWN0aW9uCiAgICBzZXF1ZW5jZXM6CiAgICAtIG5hbWU6IGFydGlmYWN0LWRlbGl2ZXJ5IAogICAgICB0cmlnZ2VyczoKICAgICAgLSBoYXJkZW5pbmcuYXJ0aWZhY3QtZGVsaXZlcnkuZmluaXNoZWQKICAgICAgdGFza3M6CiAgICAgIC0gbmFtZTogZGVwbG95bWVudAogICAgICAgIHByb3BlcnRpZXM6CiAgICAgICAgICBzdHJhdGVneTogYmx1ZV9ncmVlbgogICAgICAtIG5hbWU6IHJlbGVhc2UKICAgICAgCiAgICAtIG5hbWU6IHJlbWVkaWF0aW9uCiAgICAgIHRhc2tzOgogICAgICAtIG5hbWU6IHJlbWVkaWF0aW9uCiAgICAgIC0gbmFtZTogZXZhbHVhdGlvbg"

	secretStore := &SecretStoreMock{}
	projectRepo := &db_mock.ProjectRepoMock{}
	eventRepo := &db_mock.EventRepoMock{}
	taskSequenceRepo := &db_mock.TaskSequenceRepoMock{}
	configStore := &ConfigurationStoreMock{}

	configStore.GetProjectFunc = func(in1 string) (*keptnapimodels.Project, error) {
		return nil, nil
	}
	configStore.CreateProjectFunc = func(keptnapimodels.Project) error {
		return nil
	}

	configStore.CreateStageFunc = func(projectName string, stageName string) error {
		return nil
	}

	configStore.CreateProjectShipyardFunc = func(projectName string, resoureces []*keptnapimodels.Resource) error {
		return errors.New("whoops...")
	}

	configStore.DeleteProjectFunc = func(projectName string) error {
		return nil
	}

	secretStore.UpdateSecretFunc = func(name string, content map[string][]byte) error {
		return nil
	}

	secretStore.DeleteSecretFunc = func(name string) error {
		return nil
	}
	projectRepo.CreateProjectFunc = func(project *models.ExpandedProject) error {
		return nil
	}

	instance := NewProjectManager(configStore, secretStore, projectRepo, taskSequenceRepo, eventRepo)
	params := &operations.CreateProjectParams{
		GitRemoteURL: "git-url",
		GitToken:     "git-token",
		GitUser:      "git-user",
		Name:         stringp("my-project"),
		Shipyard:     stringp(encodedShipyard),
	}
	err, rollback := instance.Create(params)
	assert.NotNil(t, err)
	rollback()
	assert.Equal(t, "my-project", configStore.DeleteProjectCalls()[0].ProjectName)
	assert.Equal(t, "git-credentials-my-project", secretStore.DeleteSecretCalls()[0].Name)

}

func TestCreate_WhenSavingProjectInRepositoryFails_thenProjectAndSecretGetDeletedAgain(t *testing.T) {

	encodedShipyard := "YXBpVmVyc2lvbjogc3BlYy5rZXB0bi5zaC8wLjIuMApraW5kOiBTaGlweWFyZAptZXRhZGF0YToKICBuYW1lOiB0ZXN0LXNoaXB5YXJkCnNwZWM6CiAgc3RhZ2VzOgogIC0gbmFtZTogZGV2CiAgICBzZXF1ZW5jZXM6CiAgICAtIG5hbWU6IGFydGlmYWN0LWRlbGl2ZXJ5CiAgICAgIHRhc2tzOgogICAgICAtIG5hbWU6IGRlcGxveW1lbnQKICAgICAgICBwcm9wZXJ0aWVzOiAgCiAgICAgICAgICBzdHJhdGVneTogZGlyZWN0CiAgICAgIC0gbmFtZTogdGVzdAogICAgICAgIHByb3BlcnRpZXM6CiAgICAgICAgICBraW5kOiBmdW5jdGlvbmFsCiAgICAgIC0gbmFtZTogZXZhbHVhdGlvbiAKICAgICAgLSBuYW1lOiByZWxlYXNlIAoKICAtIG5hbWU6IGhhcmRlbmluZwogICAgc2VxdWVuY2VzOgogICAgLSBuYW1lOiBhcnRpZmFjdC1kZWxpdmVyeQogICAgICB0cmlnZ2VyczoKICAgICAgLSBkZXYuYXJ0aWZhY3QtZGVsaXZlcnkuZmluaXNoZWQKICAgICAgdGFza3M6CiAgICAgIC0gbmFtZTogZGVwbG95bWVudAogICAgICAgIHByb3BlcnRpZXM6IAogICAgICAgICAgc3RyYXRlZ3k6IGJsdWVfZ3JlZW5fc2VydmljZQogICAgICAtIG5hbWU6IHRlc3QKICAgICAgICBwcm9wZXJ0aWVzOiAgCiAgICAgICAgICBraW5kOiBwZXJmb3JtYW5jZQogICAgICAtIG5hbWU6IGV2YWx1YXRpb24KICAgICAgLSBuYW1lOiByZWxlYXNlCiAgICAgICAgCiAgLSBuYW1lOiBwcm9kdWN0aW9uCiAgICBzZXF1ZW5jZXM6CiAgICAtIG5hbWU6IGFydGlmYWN0LWRlbGl2ZXJ5IAogICAgICB0cmlnZ2VyczoKICAgICAgLSBoYXJkZW5pbmcuYXJ0aWZhY3QtZGVsaXZlcnkuZmluaXNoZWQKICAgICAgdGFza3M6CiAgICAgIC0gbmFtZTogZGVwbG95bWVudAogICAgICAgIHByb3BlcnRpZXM6CiAgICAgICAgICBzdHJhdGVneTogYmx1ZV9ncmVlbgogICAgICAtIG5hbWU6IHJlbGVhc2UKICAgICAgCiAgICAtIG5hbWU6IHJlbWVkaWF0aW9uCiAgICAgIHRhc2tzOgogICAgICAtIG5hbWU6IHJlbWVkaWF0aW9uCiAgICAgIC0gbmFtZTogZXZhbHVhdGlvbg"

	secretStore := &SecretStoreMock{}
	projectRepo := &db_mock.ProjectRepoMock{}
	eventRepo := &db_mock.EventRepoMock{}
	taskSequenceRepo := &db_mock.TaskSequenceRepoMock{}
	configStore := &ConfigurationStoreMock{}

	configStore.GetProjectFunc = func(in1 string) (*keptnapimodels.Project, error) { return nil, nil }
	configStore.CreateProjectFunc = func(keptnapimodels.Project) error { return nil }
	configStore.CreateStageFunc = func(projectName string, stageName string) error { return nil }
	configStore.CreateProjectShipyardFunc = func(projectName string, resoureces []*keptnapimodels.Resource) error { return nil }
	configStore.DeleteProjectFunc = func(projectName string) error { return nil }
	secretStore.UpdateSecretFunc = func(name string, content map[string][]byte) error { return nil }
	secretStore.DeleteSecretFunc = func(name string) error { return nil }
	projectRepo.CreateProjectFunc = func(project *models.ExpandedProject) error {
		return errors.New("whoops...")
	}

	instance := NewProjectManager(configStore, secretStore, projectRepo, taskSequenceRepo, eventRepo)
	params := &operations.CreateProjectParams{
		GitRemoteURL: "git-url",
		GitToken:     "git-token",
		GitUser:      "git-user",
		Name:         stringp("my-project"),
		Shipyard:     stringp(encodedShipyard),
	}
	err, rollback := instance.Create(params)
	assert.NotNil(t, err)
	rollback()
	assert.Equal(t, "my-project", configStore.DeleteProjectCalls()[0].ProjectName)
	assert.Equal(t, "git-credentials-my-project", secretStore.DeleteSecretCalls()[0].Name)

}

func TestCreate(t *testing.T) {

	encodedShipyard := "YXBpVmVyc2lvbjogc3BlYy5rZXB0bi5zaC8wLjIuMApraW5kOiBTaGlweWFyZAptZXRhZGF0YToKICBuYW1lOiB0ZXN0LXNoaXB5YXJkCnNwZWM6CiAgc3RhZ2VzOgogIC0gbmFtZTogZGV2CiAgICBzZXF1ZW5jZXM6CiAgICAtIG5hbWU6IGFydGlmYWN0LWRlbGl2ZXJ5CiAgICAgIHRhc2tzOgogICAgICAtIG5hbWU6IGRlcGxveW1lbnQKICAgICAgICBwcm9wZXJ0aWVzOiAgCiAgICAgICAgICBzdHJhdGVneTogZGlyZWN0CiAgICAgIC0gbmFtZTogdGVzdAogICAgICAgIHByb3BlcnRpZXM6CiAgICAgICAgICBraW5kOiBmdW5jdGlvbmFsCiAgICAgIC0gbmFtZTogZXZhbHVhdGlvbiAKICAgICAgLSBuYW1lOiByZWxlYXNlIAoKICAtIG5hbWU6IGhhcmRlbmluZwogICAgc2VxdWVuY2VzOgogICAgLSBuYW1lOiBhcnRpZmFjdC1kZWxpdmVyeQogICAgICB0cmlnZ2VyczoKICAgICAgLSBkZXYuYXJ0aWZhY3QtZGVsaXZlcnkuZmluaXNoZWQKICAgICAgdGFza3M6CiAgICAgIC0gbmFtZTogZGVwbG95bWVudAogICAgICAgIHByb3BlcnRpZXM6IAogICAgICAgICAgc3RyYXRlZ3k6IGJsdWVfZ3JlZW5fc2VydmljZQogICAgICAtIG5hbWU6IHRlc3QKICAgICAgICBwcm9wZXJ0aWVzOiAgCiAgICAgICAgICBraW5kOiBwZXJmb3JtYW5jZQogICAgICAtIG5hbWU6IGV2YWx1YXRpb24KICAgICAgLSBuYW1lOiByZWxlYXNlCiAgICAgICAgCiAgLSBuYW1lOiBwcm9kdWN0aW9uCiAgICBzZXF1ZW5jZXM6CiAgICAtIG5hbWU6IGFydGlmYWN0LWRlbGl2ZXJ5IAogICAgICB0cmlnZ2VyczoKICAgICAgLSBoYXJkZW5pbmcuYXJ0aWZhY3QtZGVsaXZlcnkuZmluaXNoZWQKICAgICAgdGFza3M6CiAgICAgIC0gbmFtZTogZGVwbG95bWVudAogICAgICAgIHByb3BlcnRpZXM6CiAgICAgICAgICBzdHJhdGVneTogYmx1ZV9ncmVlbgogICAgICAtIG5hbWU6IHJlbGVhc2UKICAgICAgCiAgICAtIG5hbWU6IHJlbWVkaWF0aW9uCiAgICAgIHRhc2tzOgogICAgICAtIG5hbWU6IHJlbWVkaWF0aW9uCiAgICAgIC0gbmFtZTogZXZhbHVhdGlvbg"

	secretStore := &SecretStoreMock{}
	projectRepo := &db_mock.ProjectRepoMock{}
	eventRepo := &db_mock.EventRepoMock{}
	taskSequenceRepo := &db_mock.TaskSequenceRepoMock{}
	configStore := &ConfigurationStoreMock{}

	configStore.GetProjectFunc = func(in1 string) (*keptnapimodels.Project, error) {
		return nil, nil
	}
	configStore.CreateProjectFunc = func(keptnapimodels.Project) error {
		return nil
	}

	configStore.CreateStageFunc = func(projectName string, stageName string) error {
		return nil
	}

	configStore.CreateProjectShipyardFunc = func(projectName string, resoureces []*keptnapimodels.Resource) error {
		return nil
	}

	secretStore.UpdateSecretFunc = func(name string, content map[string][]byte) error {
		return nil
	}

	projectRepo.CreateProjectFunc = func(project *models.ExpandedProject) error {
		return nil
	}

	instance := NewProjectManager(configStore, secretStore, projectRepo, taskSequenceRepo, eventRepo)
	params := &operations.CreateProjectParams{
		GitRemoteURL: "git-url",
		GitToken:     "git-token",
		GitUser:      "git-user",
		Name:         stringp("my-project"),
		Shipyard:     stringp(encodedShipyard),
	}
	instance.Create(params)
	assert.Equal(t, 3, len(configStore.CreateStageCalls()))
	assert.Equal(t, "my-project", configStore.CreateStageCalls()[0].ProjectName)
	assert.Equal(t, "dev", configStore.CreateStageCalls()[0].Stage)
	assert.Equal(t, "my-project", configStore.CreateStageCalls()[1].ProjectName)
	assert.Equal(t, "hardening", configStore.CreateStageCalls()[1].Stage)
	assert.Equal(t, "my-project", configStore.CreateStageCalls()[2].ProjectName)
	assert.Equal(t, "production", configStore.CreateStageCalls()[2].Stage)
	assert.Equal(t, "git-url", projectRepo.CreateProjectCalls()[0].Project.GitRemoteURI)
	assert.Equal(t, "git-user", projectRepo.CreateProjectCalls()[0].Project.GitUser)
	assert.Equal(t, "my-project", projectRepo.CreateProjectCalls()[0].Project.ProjectName)
	assert.Equal(t, encodedShipyard, projectRepo.CreateProjectCalls()[0].Project.Shipyard)
}

func TestUpdate_FailsWhenGettingOldSecretFails(t *testing.T) {

	secretStore := &SecretStoreMock{}
	projectRepo := &db_mock.ProjectRepoMock{}
	eventRepo := &db_mock.EventRepoMock{}
	taskSequenceRepo := &db_mock.TaskSequenceRepoMock{}
	configStore := &ConfigurationStoreMock{}

	secretStore.GetSecretFunc = func(name string) (map[string][]byte, error) {
		return nil, errors.New("Whoops...")
	}

	instance := NewProjectManager(configStore, secretStore, projectRepo, taskSequenceRepo, eventRepo)
	params := &operations.UpdateProjectParams{
		GitRemoteURL: "git-url",
		GitToken:     "git-token",
		GitUser:      "git-user",
		Name:         stringp("my-project"),
	}
	err, rollback := instance.Update(params)
	assert.NotNil(t, err)
	rollback()

}

func TestUpdate_FailsWhenGettingOldProjectFails(t *testing.T) {

	secretStore := &SecretStoreMock{}
	projectRepo := &db_mock.ProjectRepoMock{}
	eventRepo := &db_mock.EventRepoMock{}
	taskSequenceRepo := &db_mock.TaskSequenceRepoMock{}
	configStore := &ConfigurationStoreMock{}

	secretStore.GetSecretFunc = func(name string) (map[string][]byte, error) {
		return nil, nil
	}
	configStore.GetProjectFunc = func(projectName string) (*keptnapimodels.Project, error) {
		return nil, errors.New("Whoops...")
	}

	instance := NewProjectManager(configStore, secretStore, projectRepo, taskSequenceRepo, eventRepo)
	params := &operations.UpdateProjectParams{
		GitRemoteURL: "git-url",
		GitToken:     "git-token",
		GitUser:      "git-user",
		Name:         stringp("my-project"),
	}
	err, rollback := instance.Update(params)
	assert.NotNil(t, err)
	rollback()

}

func TestUpdate_FailsWhenUpdateingGitRepositorySecretFails(t *testing.T) {

	secretStore := &SecretStoreMock{}
	projectRepo := &db_mock.ProjectRepoMock{}
	eventRepo := &db_mock.EventRepoMock{}
	taskSequenceRepo := &db_mock.TaskSequenceRepoMock{}
	configStore := &ConfigurationStoreMock{}

	secretStore.GetSecretFunc = func(name string) (map[string][]byte, error) {
		return nil, nil
	}

	secretStore.UpdateSecretFunc = func(name string, content map[string][]byte) error {
		return errors.New("Whoops...")
	}
	configStore.GetProjectFunc = func(projectName string) (*keptnapimodels.Project, error) {
		return nil, nil
	}

	instance := NewProjectManager(configStore, secretStore, projectRepo, taskSequenceRepo, eventRepo)
	params := &operations.UpdateProjectParams{
		GitRemoteURL: "git-url",
		GitToken:     "git-token",
		GitUser:      "git-user",
		Name:         stringp("my-project"),
	}
	err, rollback := instance.Update(params)
	assert.NotNil(t, err)
	rollback()
	assert.Equal(t, "git-credentials-my-project", secretStore.UpdateSecretCalls()[0].Name)

}

func TestUpdate_WhenUpdateProjectInConfigurationStoreFails_ThenOldSecretGetRestored(t *testing.T) {

	secretStore := &SecretStoreMock{}
	projectRepo := &db_mock.ProjectRepoMock{}
	eventRepo := &db_mock.EventRepoMock{}
	taskSequenceRepo := &db_mock.TaskSequenceRepoMock{}
	configStore := &ConfigurationStoreMock{}

	oldSecretsEncoded, _ := json.Marshal(gitCredentials{
		User:      "my-old-user",
		Token:     "my-old-token",
		RemoteURI: "http://my-old-remote.uri",
	})

	newSecretsEncoded, _ := json.Marshal(gitCredentials{
		User:      "git-user",
		Token:     "git-token",
		RemoteURI: "git-url",
	})

	secretStore.GetSecretFunc = func(name string) (map[string][]byte, error) {

		return map[string][]byte{"git-credentials": oldSecretsEncoded}, nil
	}

	secretStore.UpdateSecretFunc = func(name string, content map[string][]byte) error {
		return nil
	}
	configStore.GetProjectFunc = func(projectName string) (*keptnapimodels.Project, error) {
		return nil, nil
	}

	configStore.UpdateProjectFunc = func(project keptnapimodels.Project) error {
		return errors.New("Whoops...")
	}

	instance := NewProjectManager(configStore, secretStore, projectRepo, taskSequenceRepo, eventRepo)
	params := &operations.UpdateProjectParams{
		GitRemoteURL: "git-url",
		GitToken:     "git-token",
		GitUser:      "git-user",
		Name:         stringp("my-project"),
	}
	err, rollback := instance.Update(params)
	assert.NotNil(t, err)
	rollback()

	expectedProjectUpdate := keptnapimodels.Project{
		GitRemoteURI: params.GitRemoteURL,
		GitToken:     params.GitToken,
		GitUser:      params.GitUser,
		ProjectName:  *params.Name,
	}
	assert.Equal(t, expectedProjectUpdate, configStore.UpdateProjectCalls()[0].Project)
	assert.Equal(t, "git-credentials-my-project", secretStore.UpdateSecretCalls()[0].Name)
	assert.Equal(t, newSecretsEncoded, secretStore.UpdateSecretCalls()[0].Content["git-credentials"])

	assert.Equal(t, "git-credentials-my-project", secretStore.UpdateSecretCalls()[1].Name)
	assert.Equal(t, oldSecretsEncoded, secretStore.UpdateSecretCalls()[1].Content["git-credentials"])
}

func TestUpdate_WhenUpdateProjectUpstreamInRepository_ThenOldProjectAndOldSecretGetRestored(t *testing.T) {

	secretStore := &SecretStoreMock{}
	projectRepo := &db_mock.ProjectRepoMock{}
	eventRepo := &db_mock.EventRepoMock{}
	taskSequenceRepo := &db_mock.TaskSequenceRepoMock{}
	configStore := &ConfigurationStoreMock{}

	oldSecretsEncoded, _ := json.Marshal(gitCredentials{
		User:      "my-old-user",
		Token:     "my-old-token",
		RemoteURI: "http://my-old-remote.uri",
	})

	newSecretsEncoded, _ := json.Marshal(gitCredentials{
		User:      "git-user",
		Token:     "git-token",
		RemoteURI: "git-url",
	})

	oldProject := keptnapimodels.Project{
		CreationDate:    "old-creationdate",
		GitRemoteURI:    "http://my-old-remote.uri",
		GitToken:        "my-old-token",
		GitUser:         "my-old-user",
		ProjectName:     "my-project",
		ShipyardVersion: "v1",
		Stages:          nil,
	}

	secretStore.GetSecretFunc = func(name string) (map[string][]byte, error) {

		return map[string][]byte{"git-credentials": oldSecretsEncoded}, nil
	}

	secretStore.UpdateSecretFunc = func(name string, content map[string][]byte) error {
		return nil
	}
	configStore.GetProjectFunc = func(projectName string) (*keptnapimodels.Project, error) {
		return &oldProject, nil
	}

	configStore.UpdateProjectFunc = func(project keptnapimodels.Project) error {
		return nil
	}

	projectRepo.UpdateProjectUpstreamFunc = func(projectName string, uri string, user string) error {
		return errors.New("Whoops...")
	}

	instance := NewProjectManager(configStore, secretStore, projectRepo, taskSequenceRepo, eventRepo)
	params := &operations.UpdateProjectParams{
		GitRemoteURL: "git-url",
		GitToken:     "git-token",
		GitUser:      "git-user",
		Name:         stringp("my-project"),
	}
	err, rollback := instance.Update(params)
	assert.NotNil(t, err)
	rollback()

	expectedProjectUpdateInConfigSvc := keptnapimodels.Project{
		GitRemoteURI: params.GitRemoteURL,
		GitToken:     params.GitToken,
		GitUser:      params.GitUser,
		ProjectName:  *params.Name,
	}

	assert.Equal(t, expectedProjectUpdateInConfigSvc, configStore.UpdateProjectCalls()[0].Project)
	assert.Equal(t, "git-credentials-my-project", secretStore.UpdateSecretCalls()[0].Name)
	assert.Equal(t, newSecretsEncoded, secretStore.UpdateSecretCalls()[0].Content["git-credentials"])
	assert.Equal(t, "my-project", projectRepo.UpdateProjectUpstreamCalls()[0].ProjectName)
	assert.Equal(t, "git-user", projectRepo.UpdateProjectUpstreamCalls()[0].User)
	assert.Equal(t, "git-url", projectRepo.UpdateProjectUpstreamCalls()[0].URI)

	assert.Equal(t, oldProject, configStore.UpdateProjectCalls()[1].Project)
	assert.Equal(t, "git-credentials-my-project", secretStore.UpdateSecretCalls()[1].Name)
	assert.Equal(t, oldSecretsEncoded, secretStore.UpdateSecretCalls()[1].Content["git-credentials"])

}

func TestUpdate(t *testing.T) {

	secretStore := &SecretStoreMock{}
	projectRepo := &db_mock.ProjectRepoMock{}
	eventRepo := &db_mock.EventRepoMock{}
	taskSequenceRepo := &db_mock.TaskSequenceRepoMock{}
	configStore := &ConfigurationStoreMock{}

	oldSecretsEncoded, _ := json.Marshal(gitCredentials{
		User:      "my-old-user",
		Token:     "my-old-token",
		RemoteURI: "http://my-old-remote.uri",
	})

	newSecretsEncoded, _ := json.Marshal(gitCredentials{
		User:      "git-user",
		Token:     "git-token",
		RemoteURI: "git-url",
	})

	oldProject := keptnapimodels.Project{
		CreationDate:    "old-creationdate",
		GitRemoteURI:    "http://my-old-remote.uri",
		GitToken:        "my-old-token",
		GitUser:         "my-old-user",
		ProjectName:     "my-project",
		ShipyardVersion: "v1",
		Stages:          nil,
	}

	secretStore.GetSecretFunc = func(name string) (map[string][]byte, error) {

		return map[string][]byte{"git-credentials": oldSecretsEncoded}, nil
	}

	secretStore.UpdateSecretFunc = func(name string, content map[string][]byte) error {
		return nil
	}
	configStore.GetProjectFunc = func(projectName string) (*keptnapimodels.Project, error) {
		return &oldProject, nil
	}

	configStore.UpdateProjectFunc = func(project keptnapimodels.Project) error {
		return nil
	}

	projectRepo.UpdateProjectUpstreamFunc = func(projectName string, uri string, user string) error {
		return nil
	}

	instance := NewProjectManager(configStore, secretStore, projectRepo, taskSequenceRepo, eventRepo)
	params := &operations.UpdateProjectParams{
		GitRemoteURL: "git-url",
		GitToken:     "git-token",
		GitUser:      "git-user",
		Name:         stringp("my-project"),
	}
	err, rollback := instance.Update(params)
	assert.Nil(t, err)
	rollback()

	expectedProjectUpdateInConfigSvc := keptnapimodels.Project{
		GitRemoteURI: params.GitRemoteURL,
		GitToken:     params.GitToken,
		GitUser:      params.GitUser,
		ProjectName:  *params.Name,
	}

	assert.Equal(t, expectedProjectUpdateInConfigSvc, configStore.UpdateProjectCalls()[0].Project)
	assert.Equal(t, "git-credentials-my-project", secretStore.UpdateSecretCalls()[0].Name)
	assert.Equal(t, newSecretsEncoded, secretStore.UpdateSecretCalls()[0].Content["git-credentials"])
	assert.Equal(t, "my-project", projectRepo.UpdateProjectUpstreamCalls()[0].ProjectName)
	assert.Equal(t, "git-user", projectRepo.UpdateProjectUpstreamCalls()[0].User)
	assert.Equal(t, "git-url", projectRepo.UpdateProjectUpstreamCalls()[0].URI)

}

func TestDelete(t *testing.T) {

	secretStore := &SecretStoreMock{}
	projectRepo := &db_mock.ProjectRepoMock{}
	eventRepo := &db_mock.EventRepoMock{}
	taskSequenceRepo := &db_mock.TaskSequenceRepoMock{}
	configStore := &ConfigurationStoreMock{}

	configStore.GetProjectFunc = func(projectName string) (*keptnapimodels.Project, error) {
		p := keptnapimodels.Project{
			CreationDate:    "creationdate",
			GitRemoteURI:    "http://my-remote.uri",
			GitToken:        "my-token",
			GitUser:         "my-user",
			ProjectName:     "my-project",
			ShipyardVersion: "v1",
			Stages:          nil,
		}
		return &p, nil
	}

	secretEncoded, _ := json.Marshal(gitCredentials{
		User:      "my-user",
		Token:     "my-token",
		RemoteURI: "http://my-remote.uri",
	})

	secretStore.GetSecretFunc = func(name string) (map[string][]byte, error) {

		return map[string][]byte{"git-credentials": secretEncoded}, nil
	}

	secretStore.DeleteSecretFunc = func(name string) error {
		return nil
	}

	configStore.DeleteProjectFunc = func(projectName string) error {
		return nil
	}

	configStore.GetProjectResourceFunc = func(projectName string, resourceURI string) (*keptnapimodels.Resource, error) {
		resource := keptnapimodels.Resource{}
		return &resource, nil
	}
	eventRepo.DeleteEventCollectionsFunc = func(project string) error {
		return nil
	}

	taskSequenceRepo.DeleteTaskSequenceCollectionFunc = func(project string) error {
		return nil
	}

	projectRepo.DeleteProjectFunc = func(projectName string) error {
		return nil
	}

	instance := NewProjectManager(configStore, secretStore, projectRepo, taskSequenceRepo, eventRepo)
	instance.Delete("my-project")
}
