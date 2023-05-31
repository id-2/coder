package dbmetrics

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/exp/slices"

	"github.com/coder/coder/coderd/database"
	"github.com/coder/coder/coderd/rbac"
)

const wrapname = "dbmetrics.metricsStore"

// New returns a database.Store that registers metrics for all queries to reg.
func New(s database.Store, reg prometheus.Registerer) database.Store {
	// Don't double-wrap.
	if slices.Contains(s.Wrappers(), wrapname) {
		return s
	}
	queryLatencies := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "coderd",
		Subsystem: "db",
		Name:      "query_latencies_seconds",
		Help:      "Latency distribution of queries in seconds.",
		Buckets:   prometheus.DefBuckets,
	}, []string{"query"})
	txDuration := prometheus.NewHistogram(prometheus.HistogramOpts{
		Namespace: "coderd",
		Subsystem: "db",
		Name:      "tx_duration_seconds",
		Help:      "Duration of transactions in seconds.",
		Buckets:   prometheus.DefBuckets,
	})
	reg.MustRegister(queryLatencies)
	reg.MustRegister(txDuration)
	return &metricsStore{
		s:              s,
		queryLatencies: queryLatencies,
		txDuration:     txDuration,
	}
}

var _ database.Store = (*metricsStore)(nil)

type metricsStore struct {
	s              database.Store
	queryLatencies *prometheus.HistogramVec
	txDuration     prometheus.Histogram
}

func (m metricsStore) Wrappers() []string {
	return append(m.s.Wrappers(), wrapname)
}

func (m metricsStore) Ping(ctx context.Context) (time.Duration, error) {
	start := time.Now()
	duration, err := m.s.Ping(ctx)
	m.queryLatencies.WithLabelValues("Ping").Observe(time.Since(start).Seconds())
	return duration, err
}

func (m metricsStore) InTx(f func(database.Store) error, options *sql.TxOptions) error {
	start := time.Now()
	err := m.s.InTx(f, options)
	m.txDuration.Observe(time.Since(start).Seconds())
	return err
}

func (m metricsStore) AcquireLock(ctx context.Context, pgAdvisoryXactLock int64) error {
	start := time.Now()
	err := m.s.AcquireLock(ctx, pgAdvisoryXactLock)
	m.queryLatencies.WithLabelValues("AcquireLock").Observe(time.Since(start).Seconds())
	return err
}

func (m metricsStore) AcquireProvisionerJob(ctx context.Context, arg database.AcquireProvisionerJobParams) (database.ProvisionerJob, error) {
	start := time.Now()
	provisionerJob, err := m.s.AcquireProvisionerJob(ctx, arg)
	m.queryLatencies.WithLabelValues("AcquireProvisionerJob").Observe(time.Since(start).Seconds())
	return provisionerJob, err
}

func (m metricsStore) DeleteAPIKeyByID(ctx context.Context, id string) error {
	start := time.Now()
	err := m.s.DeleteAPIKeyByID(ctx, id)
	m.queryLatencies.WithLabelValues("DeleteAPIKeyByID").Observe(time.Since(start).Seconds())
	return err
}

func (m metricsStore) DeleteAPIKeysByUserID(ctx context.Context, userID uuid.UUID) error {
	start := time.Now()
	err := m.s.DeleteAPIKeysByUserID(ctx, userID)
	m.queryLatencies.WithLabelValues("DeleteAPIKeysByUserID").Observe(time.Since(start).Seconds())
	return err
}

func (m metricsStore) DeleteApplicationConnectAPIKeysByUserID(ctx context.Context, userID uuid.UUID) error {
	start := time.Now()
	err := m.s.DeleteApplicationConnectAPIKeysByUserID(ctx, userID)
	m.queryLatencies.WithLabelValues("DeleteApplicationConnectAPIKeysByUserID").Observe(time.Since(start).Seconds())
	return err
}

func (m metricsStore) DeleteGitSSHKey(ctx context.Context, userID uuid.UUID) error {
	start := time.Now()
	err := m.s.DeleteGitSSHKey(ctx, userID)
	m.queryLatencies.WithLabelValues("DeleteGitSSHKey").Observe(time.Since(start).Seconds())
	return err
}

func (m metricsStore) DeleteGroupByID(ctx context.Context, id uuid.UUID) error {
	start := time.Now()
	err := m.s.DeleteGroupByID(ctx, id)
	m.queryLatencies.WithLabelValues("DeleteGroupByID").Observe(time.Since(start).Seconds())
	return err
}

func (m metricsStore) DeleteGroupMemberFromGroup(ctx context.Context, arg database.DeleteGroupMemberFromGroupParams) error {
	start := time.Now()
	err := m.s.DeleteGroupMemberFromGroup(ctx, arg)
	m.queryLatencies.WithLabelValues("DeleteGroupMemberFromGroup").Observe(time.Since(start).Seconds())
	return err
}

func (m metricsStore) DeleteGroupMembersByOrgAndUser(ctx context.Context, arg database.DeleteGroupMembersByOrgAndUserParams) error {
	start := time.Now()
	err := m.s.DeleteGroupMembersByOrgAndUser(ctx, arg)
	m.queryLatencies.WithLabelValues("DeleteGroupMembersByOrgAndUser").Observe(time.Since(start).Seconds())
	return err
}

func (m metricsStore) DeleteLicense(ctx context.Context, id int32) (int32, error) {
	start := time.Now()
	licenseID, err := m.s.DeleteLicense(ctx, id)
	m.queryLatencies.WithLabelValues("DeleteLicense").Observe(time.Since(start).Seconds())
	return licenseID, err
}

func (m metricsStore) DeleteOldWorkspaceAgentStartupLogs(ctx context.Context) error {
	start := time.Now()
	err := m.s.DeleteOldWorkspaceAgentStartupLogs(ctx)
	m.queryLatencies.WithLabelValues("DeleteOldWorkspaceAgentStartupLogs").Observe(time.Since(start).Seconds())
	return err
}

func (m metricsStore) DeleteOldWorkspaceAgentStats(ctx context.Context) error {
	start := time.Now()
	err := m.s.DeleteOldWorkspaceAgentStats(ctx)
	m.queryLatencies.WithLabelValues("DeleteOldWorkspaceAgentStats").Observe(time.Since(start).Seconds())
	return err
}

func (m metricsStore) DeleteParameterValueByID(ctx context.Context, id uuid.UUID) error {
	start := time.Now()
	err := m.s.DeleteParameterValueByID(ctx, id)
	m.queryLatencies.WithLabelValues("DeleteParameterValueByID").Observe(time.Since(start).Seconds())
	return err
}

func (m metricsStore) DeleteReplicasUpdatedBefore(ctx context.Context, updatedAt time.Time) error {
	start := time.Now()
	err := m.s.DeleteReplicasUpdatedBefore(ctx, updatedAt)
	m.queryLatencies.WithLabelValues("DeleteReplicasUpdatedBefore").Observe(time.Since(start).Seconds())
	return err
}

func (m metricsStore) GetAPIKeyByID(ctx context.Context, id string) (database.APIKey, error) {
	start := time.Now()
	apiKey, err := m.s.GetAPIKeyByID(ctx, id)
	m.queryLatencies.WithLabelValues("GetAPIKeyByID").Observe(time.Since(start).Seconds())
	return apiKey, err
}

func (m metricsStore) GetAPIKeyByName(ctx context.Context, arg database.GetAPIKeyByNameParams) (database.APIKey, error) {
	start := time.Now()
	apiKey, err := m.s.GetAPIKeyByName(ctx, arg)
	m.queryLatencies.WithLabelValues("GetAPIKeyByName").Observe(time.Since(start).Seconds())
	return apiKey, err
}

func (m metricsStore) GetAPIKeysByLoginType(ctx context.Context, loginType database.LoginType) ([]database.APIKey, error) {
	start := time.Now()
	apiKeys, err := m.s.GetAPIKeysByLoginType(ctx, loginType)
	m.queryLatencies.WithLabelValues("GetAPIKeysByLoginType").Observe(time.Since(start).Seconds())
	return apiKeys, err
}

func (m metricsStore) GetAPIKeysByUserID(ctx context.Context, arg database.GetAPIKeysByUserIDParams) ([]database.APIKey, error) {
	start := time.Now()
	apiKeys, err := m.s.GetAPIKeysByUserID(ctx, arg)
	m.queryLatencies.WithLabelValues("GetAPIKeysByUserID").Observe(time.Since(start).Seconds())
	return apiKeys, err
}

func (m metricsStore) GetAPIKeysLastUsedAfter(ctx context.Context, lastUsed time.Time) ([]database.APIKey, error) {
	start := time.Now()
	apiKeys, err := m.s.GetAPIKeysLastUsedAfter(ctx, lastUsed)
	m.queryLatencies.WithLabelValues("GetAPIKeysLastUsedAfter").Observe(time.Since(start).Seconds())
	return apiKeys, err
}

func (m metricsStore) GetActiveUserCount(ctx context.Context) (int64, error) {
	start := time.Now()
	count, err := m.s.GetActiveUserCount(ctx)
	m.queryLatencies.WithLabelValues("GetActiveUserCount").Observe(time.Since(start).Seconds())
	return count, err
}

func (m metricsStore) GetAppSecurityKey(ctx context.Context) (string, error) {
	start := time.Now()
	key, err := m.s.GetAppSecurityKey(ctx)
	m.queryLatencies.WithLabelValues("GetAppSecurityKey").Observe(time.Since(start).Seconds())
	return key, err
}

func (m metricsStore) GetAuditLogsOffset(ctx context.Context, arg database.GetAuditLogsOffsetParams) ([]database.GetAuditLogsOffsetRow, error) {
	start := time.Now()
	rows, err := m.s.GetAuditLogsOffset(ctx, arg)
	m.queryLatencies.WithLabelValues("GetAuditLogsOffset").Observe(time.Since(start).Seconds())
	return rows, err
}

func (m metricsStore) GetAuthorizationUserRoles(ctx context.Context, userID uuid.UUID) (database.GetAuthorizationUserRolesRow, error) {
	start := time.Now()
	row, err := m.s.GetAuthorizationUserRoles(ctx, userID)
	m.queryLatencies.WithLabelValues("GetAuthorizationUserRoles").Observe(time.Since(start).Seconds())
	return row, err
}

func (m metricsStore) GetDERPMeshKey(ctx context.Context) (string, error) {
	start := time.Now()
	key, err := m.s.GetDERPMeshKey(ctx)
	m.queryLatencies.WithLabelValues("GetDERPMeshKey").Observe(time.Since(start).Seconds())
	return key, err
}

func (m metricsStore) GetDeploymentDAUs(ctx context.Context, tzOffset int32) ([]database.GetDeploymentDAUsRow, error) {
	start := time.Now()
	rows, err := m.s.GetDeploymentDAUs(ctx, tzOffset)
	m.queryLatencies.WithLabelValues("GetDeploymentDAUs").Observe(time.Since(start).Seconds())
	return rows, err
}

func (m metricsStore) GetDeploymentID(ctx context.Context) (string, error) {
	start := time.Now()
	id, err := m.s.GetDeploymentID(ctx)
	m.queryLatencies.WithLabelValues("GetDeploymentID").Observe(time.Since(start).Seconds())
	return id, err
}

func (m metricsStore) GetDeploymentWorkspaceAgentStats(ctx context.Context, createdAt time.Time) (database.GetDeploymentWorkspaceAgentStatsRow, error) {
	start := time.Now()
	row, err := m.s.GetDeploymentWorkspaceAgentStats(ctx, createdAt)
	m.queryLatencies.WithLabelValues("GetDeploymentWorkspaceAgentStats").Observe(time.Since(start).Seconds())
	return row, err
}

func (m metricsStore) GetDeploymentWorkspaceStats(ctx context.Context) (database.GetDeploymentWorkspaceStatsRow, error) {
	start := time.Now()
	row, err := m.s.GetDeploymentWorkspaceStats(ctx)
	m.queryLatencies.WithLabelValues("GetDeploymentWorkspaceStats").Observe(time.Since(start).Seconds())
	return row, err
}

func (m metricsStore) GetFileByHashAndCreator(ctx context.Context, arg database.GetFileByHashAndCreatorParams) (database.File, error) {
	start := time.Now()
	file, err := m.s.GetFileByHashAndCreator(ctx, arg)
	m.queryLatencies.WithLabelValues("GetFileByHashAndCreator").Observe(time.Since(start).Seconds())
	return file, err
}

func (m metricsStore) GetFileByID(ctx context.Context, id uuid.UUID) (database.File, error) {
	start := time.Now()
	file, err := m.s.GetFileByID(ctx, id)
	m.queryLatencies.WithLabelValues("GetFileByID").Observe(time.Since(start).Seconds())
	return file, err
}

func (m metricsStore) GetFileTemplates(ctx context.Context, fileID uuid.UUID) ([]database.GetFileTemplatesRow, error) {
	start := time.Now()
	rows, err := m.s.GetFileTemplates(ctx, fileID)
	m.queryLatencies.WithLabelValues("GetFileTemplates").Observe(time.Since(start).Seconds())
	return rows, err
}

func (m metricsStore) GetFilteredUserCount(ctx context.Context, arg database.GetFilteredUserCountParams) (int64, error) {
	start := time.Now()
	count, err := m.s.GetFilteredUserCount(ctx, arg)
	m.queryLatencies.WithLabelValues("GetFilteredUserCount").Observe(time.Since(start).Seconds())
	return count, err
}

func (m metricsStore) GetGitAuthLink(ctx context.Context, arg database.GetGitAuthLinkParams) (database.GitAuthLink, error) {
	start := time.Now()
	link, err := m.s.GetGitAuthLink(ctx, arg)
	m.queryLatencies.WithLabelValues("GetGitAuthLink").Observe(time.Since(start).Seconds())
	return link, err
}

func (m metricsStore) GetGitSSHKey(ctx context.Context, userID uuid.UUID) (database.GitSSHKey, error) {
	start := time.Now()
	key, err := m.s.GetGitSSHKey(ctx, userID)
	m.queryLatencies.WithLabelValues("GetGitSSHKey").Observe(time.Since(start).Seconds())
	return key, err
}

func (m metricsStore) GetGroupByID(ctx context.Context, id uuid.UUID) (database.Group, error) {
	start := time.Now()
	group, err := m.s.GetGroupByID(ctx, id)
	m.queryLatencies.WithLabelValues("GetGroupByID").Observe(time.Since(start).Seconds())
	return group, err
}

func (m metricsStore) GetGroupByOrgAndName(ctx context.Context, arg database.GetGroupByOrgAndNameParams) (database.Group, error) {
	start := time.Now()
	group, err := m.s.GetGroupByOrgAndName(ctx, arg)
	m.queryLatencies.WithLabelValues("GetGroupByOrgAndName").Observe(time.Since(start).Seconds())
	return group, err
}

func (m metricsStore) GetGroupMembers(ctx context.Context, groupID uuid.UUID) ([]database.User, error) {
	start := time.Now()
	users, err := m.s.GetGroupMembers(ctx, groupID)
	m.queryLatencies.WithLabelValues("GetGroupMembers").Observe(time.Since(start).Seconds())
	return users, err
}

func (m metricsStore) GetGroupsByOrganizationID(ctx context.Context, organizationID uuid.UUID) ([]database.Group, error) {
	start := time.Now()
	groups, err := m.s.GetGroupsByOrganizationID(ctx, organizationID)
	m.queryLatencies.WithLabelValues("GetGroupsByOrganizationID").Observe(time.Since(start).Seconds())
	return groups, err
}

func (m metricsStore) GetLastUpdateCheck(ctx context.Context) (string, error) {
	start := time.Now()
	version, err := m.s.GetLastUpdateCheck(ctx)
	m.queryLatencies.WithLabelValues("GetLastUpdateCheck").Observe(time.Since(start).Seconds())
	return version, err
}

func (m metricsStore) GetLatestWorkspaceBuildByWorkspaceID(ctx context.Context, workspaceID uuid.UUID) (database.WorkspaceBuild, error) {
	start := time.Now()
	build, err := m.s.GetLatestWorkspaceBuildByWorkspaceID(ctx, workspaceID)
	m.queryLatencies.WithLabelValues("GetLatestWorkspaceBuildByWorkspaceID").Observe(time.Since(start).Seconds())
	return build, err
}

func (m metricsStore) GetLatestWorkspaceBuilds(ctx context.Context) ([]database.WorkspaceBuild, error) {
	start := time.Now()
	builds, err := m.s.GetLatestWorkspaceBuilds(ctx)
	m.queryLatencies.WithLabelValues("GetLatestWorkspaceBuilds").Observe(time.Since(start).Seconds())
	return builds, err
}

func (m metricsStore) GetLatestWorkspaceBuildsByWorkspaceIDs(ctx context.Context, ids []uuid.UUID) ([]database.WorkspaceBuild, error) {
	start := time.Now()
	builds, err := m.s.GetLatestWorkspaceBuildsByWorkspaceIDs(ctx, ids)
	m.queryLatencies.WithLabelValues("GetLatestWorkspaceBuildsByWorkspaceIDs").Observe(time.Since(start).Seconds())
	return builds, err
}

func (m metricsStore) GetLicenseByID(ctx context.Context, id int32) (database.License, error) {
	start := time.Now()
	license, err := m.s.GetLicenseByID(ctx, id)
	m.queryLatencies.WithLabelValues("GetLicenseByID").Observe(time.Since(start).Seconds())
	return license, err
}

func (m metricsStore) GetLicenses(ctx context.Context) ([]database.License, error) {
	start := time.Now()
	licenses, err := m.s.GetLicenses(ctx)
	m.queryLatencies.WithLabelValues("GetLicenses").Observe(time.Since(start).Seconds())
	return licenses, err
}

func (m metricsStore) GetLogoURL(ctx context.Context) (string, error) {
	start := time.Now()
	url, err := m.s.GetLogoURL(ctx)
	m.queryLatencies.WithLabelValues("GetLogoURL").Observe(time.Since(start).Seconds())
	return url, err
}

func (m metricsStore) GetOrganizationByID(ctx context.Context, id uuid.UUID) (database.Organization, error) {
	start := time.Now()
	organization, err := m.s.GetOrganizationByID(ctx, id)
	m.queryLatencies.WithLabelValues("GetOrganizationByID").Observe(time.Since(start).Seconds())
	return organization, err
}

func (m metricsStore) GetOrganizationByName(ctx context.Context, name string) (database.Organization, error) {
	start := time.Now()
	organization, err := m.s.GetOrganizationByName(ctx, name)
	m.queryLatencies.WithLabelValues("GetOrganizationByName").Observe(time.Since(start).Seconds())
	return organization, err
}

func (m metricsStore) GetOrganizationIDsByMemberIDs(ctx context.Context, ids []uuid.UUID) ([]database.GetOrganizationIDsByMemberIDsRow, error) {
	start := time.Now()
	organizations, err := m.s.GetOrganizationIDsByMemberIDs(ctx, ids)
	m.queryLatencies.WithLabelValues("GetOrganizationIDsByMemberIDs").Observe(time.Since(start).Seconds())
	return organizations, err
}

func (m metricsStore) GetOrganizationMemberByUserID(ctx context.Context, arg database.GetOrganizationMemberByUserIDParams) (database.OrganizationMember, error) {
	start := time.Now()
	member, err := m.s.GetOrganizationMemberByUserID(ctx, arg)
	m.queryLatencies.WithLabelValues("GetOrganizationMemberByUserID").Observe(time.Since(start).Seconds())
	return member, err
}

func (m metricsStore) GetOrganizationMembershipsByUserID(ctx context.Context, userID uuid.UUID) ([]database.OrganizationMember, error) {
	start := time.Now()
	memberships, err := m.s.GetOrganizationMembershipsByUserID(ctx, userID)
	m.queryLatencies.WithLabelValues("GetOrganizationMembershipsByUserID").Observe(time.Since(start).Seconds())
	return memberships, err
}

func (m metricsStore) GetOrganizations(ctx context.Context) ([]database.Organization, error) {
	start := time.Now()
	organizations, err := m.s.GetOrganizations(ctx)
	m.queryLatencies.WithLabelValues("GetOrganizations").Observe(time.Since(start).Seconds())
	return organizations, err
}

func (m metricsStore) GetOrganizationsByUserID(ctx context.Context, userID uuid.UUID) ([]database.Organization, error) {
	start := time.Now()
	organizations, err := m.s.GetOrganizationsByUserID(ctx, userID)
	m.queryLatencies.WithLabelValues("GetOrganizationsByUserID").Observe(time.Since(start).Seconds())
	return organizations, err
}

func (m metricsStore) GetParameterSchemasByJobID(ctx context.Context, jobID uuid.UUID) ([]database.ParameterSchema, error) {
	start := time.Now()
	schemas, err := m.s.GetParameterSchemasByJobID(ctx, jobID)
	m.queryLatencies.WithLabelValues("GetParameterSchemasByJobID").Observe(time.Since(start).Seconds())
	return schemas, err
}

func (m metricsStore) GetParameterSchemasCreatedAfter(ctx context.Context, createdAt time.Time) ([]database.ParameterSchema, error) {
	start := time.Now()
	schemas, err := m.s.GetParameterSchemasCreatedAfter(ctx, createdAt)
	m.queryLatencies.WithLabelValues("GetParameterSchemasCreatedAfter").Observe(time.Since(start).Seconds())
	return schemas, err
}

func (m metricsStore) GetParameterValueByScopeAndName(ctx context.Context, arg database.GetParameterValueByScopeAndNameParams) (database.ParameterValue, error) {
	start := time.Now()
	value, err := m.s.GetParameterValueByScopeAndName(ctx, arg)
	m.queryLatencies.WithLabelValues("GetParameterValueByScopeAndName").Observe(time.Since(start).Seconds())
	return value, err
}

func (m metricsStore) GetPreviousTemplateVersion(ctx context.Context, arg database.GetPreviousTemplateVersionParams) (database.TemplateVersion, error) {
	start := time.Now()
	version, err := m.s.GetPreviousTemplateVersion(ctx, arg)
	m.queryLatencies.WithLabelValues("GetPreviousTemplateVersion").Observe(time.Since(start).Seconds())
	return version, err
}

func (m metricsStore) GetProvisionerDaemons(ctx context.Context) ([]database.ProvisionerDaemon, error) {
	start := time.Now()
	daemons, err := m.s.GetProvisionerDaemons(ctx)
	m.queryLatencies.WithLabelValues("GetProvisionerDaemons").Observe(time.Since(start).Seconds())
	return daemons, err
}

func (m metricsStore) GetProvisionerJobByID(ctx context.Context, id uuid.UUID) (database.ProvisionerJob, error) {
	start := time.Now()
	job, err := m.s.GetProvisionerJobByID(ctx, id)
	m.queryLatencies.WithLabelValues("GetProvisionerJobByID").Observe(time.Since(start).Seconds())
	return job, err
}

func (m metricsStore) GetProvisionerJobsByIDs(ctx context.Context, ids []uuid.UUID) ([]database.ProvisionerJob, error) {
	start := time.Now()
	jobs, err := m.s.GetProvisionerJobsByIDs(ctx, ids)
	m.queryLatencies.WithLabelValues("GetProvisionerJobsByIDs").Observe(time.Since(start).Seconds())
	return jobs, err
}

func (m metricsStore) GetProvisionerJobsCreatedAfter(ctx context.Context, createdAt time.Time) ([]database.ProvisionerJob, error) {
	start := time.Now()
	jobs, err := m.s.GetProvisionerJobsCreatedAfter(ctx, createdAt)
	m.queryLatencies.WithLabelValues("GetProvisionerJobsCreatedAfter").Observe(time.Since(start).Seconds())
	return jobs, err
}

func (m metricsStore) GetProvisionerLogsAfterID(ctx context.Context, arg database.GetProvisionerLogsAfterIDParams) ([]database.ProvisionerJobLog, error) {
	start := time.Now()
	logs, err := m.s.GetProvisionerLogsAfterID(ctx, arg)
	m.queryLatencies.WithLabelValues("GetProvisionerLogsAfterID").Observe(time.Since(start).Seconds())
	return logs, err
}

func (m metricsStore) GetQuotaAllowanceForUser(ctx context.Context, userID uuid.UUID) (int64, error) {
	start := time.Now()
	allowance, err := m.s.GetQuotaAllowanceForUser(ctx, userID)
	m.queryLatencies.WithLabelValues("GetQuotaAllowanceForUser").Observe(time.Since(start).Seconds())
	return allowance, err
}

func (m metricsStore) GetQuotaConsumedForUser(ctx context.Context, ownerID uuid.UUID) (int64, error) {
	start := time.Now()
	consumed, err := m.s.GetQuotaConsumedForUser(ctx, ownerID)
	m.queryLatencies.WithLabelValues("GetQuotaConsumedForUser").Observe(time.Since(start).Seconds())
	return consumed, err
}

func (m metricsStore) GetReplicasUpdatedAfter(ctx context.Context, updatedAt time.Time) ([]database.Replica, error) {
	start := time.Now()
	replicas, err := m.s.GetReplicasUpdatedAfter(ctx, updatedAt)
	m.queryLatencies.WithLabelValues("GetReplicasUpdatedAfter").Observe(time.Since(start).Seconds())
	return replicas, err
}

func (m metricsStore) GetServiceBanner(ctx context.Context) (string, error) {
	start := time.Now()
	banner, err := m.s.GetServiceBanner(ctx)
	m.queryLatencies.WithLabelValues("GetServiceBanner").Observe(time.Since(start).Seconds())
	return banner, err
}

func (m metricsStore) GetTemplateAverageBuildTime(ctx context.Context, arg database.GetTemplateAverageBuildTimeParams) (database.GetTemplateAverageBuildTimeRow, error) {
	start := time.Now()
	buildTime, err := m.s.GetTemplateAverageBuildTime(ctx, arg)
	m.queryLatencies.WithLabelValues("GetTemplateAverageBuildTime").Observe(time.Since(start).Seconds())
	return buildTime, err
}

func (m metricsStore) GetTemplateByID(ctx context.Context, id uuid.UUID) (database.Template, error) {
	start := time.Now()
	template, err := m.s.GetTemplateByID(ctx, id)
	m.queryLatencies.WithLabelValues("GetTemplateByID").Observe(time.Since(start).Seconds())
	return template, err
}

func (m metricsStore) GetTemplateByOrganizationAndName(ctx context.Context, arg database.GetTemplateByOrganizationAndNameParams) (database.Template, error) {
	start := time.Now()
	template, err := m.s.GetTemplateByOrganizationAndName(ctx, arg)
	m.queryLatencies.WithLabelValues("GetTemplateByOrganizationAndName").Observe(time.Since(start).Seconds())
	return template, err
}

func (m metricsStore) GetTemplateDAUs(ctx context.Context, arg database.GetTemplateDAUsParams) ([]database.GetTemplateDAUsRow, error) {
	start := time.Now()
	daus, err := m.s.GetTemplateDAUs(ctx, arg)
	m.queryLatencies.WithLabelValues("GetTemplateDAUs").Observe(time.Since(start).Seconds())
	return daus, err
}

func (m metricsStore) GetTemplateVersionByID(ctx context.Context, id uuid.UUID) (database.TemplateVersion, error) {
	start := time.Now()
	version, err := m.s.GetTemplateVersionByID(ctx, id)
	m.queryLatencies.WithLabelValues("GetTemplateVersionByID").Observe(time.Since(start).Seconds())
	return version, err
}

func (m metricsStore) GetTemplateVersionByJobID(ctx context.Context, jobID uuid.UUID) (database.TemplateVersion, error) {
	start := time.Now()
	version, err := m.s.GetTemplateVersionByJobID(ctx, jobID)
	m.queryLatencies.WithLabelValues("GetTemplateVersionByJobID").Observe(time.Since(start).Seconds())
	return version, err
}

func (m metricsStore) GetTemplateVersionByTemplateIDAndName(ctx context.Context, arg database.GetTemplateVersionByTemplateIDAndNameParams) (database.TemplateVersion, error) {
	start := time.Now()
	version, err := m.s.GetTemplateVersionByTemplateIDAndName(ctx, arg)
	m.queryLatencies.WithLabelValues("GetTemplateVersionByTemplateIDAndName").Observe(time.Since(start).Seconds())
	return version, err
}

func (m metricsStore) GetTemplateVersionParameters(ctx context.Context, templateVersionID uuid.UUID) ([]database.TemplateVersionParameter, error) {
	start := time.Now()
	parameters, err := m.s.GetTemplateVersionParameters(ctx, templateVersionID)
	m.queryLatencies.WithLabelValues("GetTemplateVersionParameters").Observe(time.Since(start).Seconds())
	return parameters, err
}

func (m metricsStore) GetTemplateVersionVariables(ctx context.Context, templateVersionID uuid.UUID) ([]database.TemplateVersionVariable, error) {
	start := time.Now()
	variables, err := m.s.GetTemplateVersionVariables(ctx, templateVersionID)
	m.queryLatencies.WithLabelValues("GetTemplateVersionVariables").Observe(time.Since(start).Seconds())
	return variables, err
}

func (m metricsStore) GetTemplateVersionsByIDs(ctx context.Context, ids []uuid.UUID) ([]database.TemplateVersion, error) {
	start := time.Now()
	versions, err := m.s.GetTemplateVersionsByIDs(ctx, ids)
	m.queryLatencies.WithLabelValues("GetTemplateVersionsByIDs").Observe(time.Since(start).Seconds())
	return versions, err
}

func (m metricsStore) GetTemplateVersionsByTemplateID(ctx context.Context, arg database.GetTemplateVersionsByTemplateIDParams) ([]database.TemplateVersion, error) {
	start := time.Now()
	versions, err := m.s.GetTemplateVersionsByTemplateID(ctx, arg)
	m.queryLatencies.WithLabelValues("GetTemplateVersionsByTemplateID").Observe(time.Since(start).Seconds())
	return versions, err
}

func (m metricsStore) GetTemplateVersionsCreatedAfter(ctx context.Context, createdAt time.Time) ([]database.TemplateVersion, error) {
	start := time.Now()
	versions, err := m.s.GetTemplateVersionsCreatedAfter(ctx, createdAt)
	m.queryLatencies.WithLabelValues("GetTemplateVersionsCreatedAfter").Observe(time.Since(start).Seconds())
	return versions, err
}

func (m metricsStore) GetTemplates(ctx context.Context) ([]database.Template, error) {
	start := time.Now()
	templates, err := m.s.GetTemplates(ctx)
	m.queryLatencies.WithLabelValues("GetTemplates").Observe(time.Since(start).Seconds())
	return templates, err
}

func (m metricsStore) GetTemplatesWithFilter(ctx context.Context, arg database.GetTemplatesWithFilterParams) ([]database.Template, error) {
	start := time.Now()
	templates, err := m.s.GetTemplatesWithFilter(ctx, arg)
	m.queryLatencies.WithLabelValues("GetTemplatesWithFilter").Observe(time.Since(start).Seconds())
	return templates, err
}

func (m metricsStore) GetUnexpiredLicenses(ctx context.Context) ([]database.License, error) {
	start := time.Now()
	licenses, err := m.s.GetUnexpiredLicenses(ctx)
	m.queryLatencies.WithLabelValues("GetUnexpiredLicenses").Observe(time.Since(start).Seconds())
	return licenses, err
}

func (m metricsStore) GetUserByEmailOrUsername(ctx context.Context, arg database.GetUserByEmailOrUsernameParams) (database.User, error) {
	start := time.Now()
	user, err := m.s.GetUserByEmailOrUsername(ctx, arg)
	m.queryLatencies.WithLabelValues("GetUserByEmailOrUsername").Observe(time.Since(start).Seconds())
	return user, err
}

func (m metricsStore) GetUserByID(ctx context.Context, id uuid.UUID) (database.User, error) {
	start := time.Now()
	user, err := m.s.GetUserByID(ctx, id)
	m.queryLatencies.WithLabelValues("GetUserByID").Observe(time.Since(start).Seconds())
	return user, err
}

func (m metricsStore) GetUserCount(ctx context.Context) (int64, error) {
	start := time.Now()
	count, err := m.s.GetUserCount(ctx)
	m.queryLatencies.WithLabelValues("GetUserCount").Observe(time.Since(start).Seconds())
	return count, err
}

func (m metricsStore) GetUserLinkByLinkedID(ctx context.Context, linkedID string) (database.UserLink, error) {
	start := time.Now()
	link, err := m.s.GetUserLinkByLinkedID(ctx, linkedID)
	m.queryLatencies.WithLabelValues("GetUserLinkByLinkedID").Observe(time.Since(start).Seconds())
	return link, err
}

func (m metricsStore) GetUserLinkByUserIDLoginType(ctx context.Context, arg database.GetUserLinkByUserIDLoginTypeParams) (database.UserLink, error) {
	start := time.Now()
	link, err := m.s.GetUserLinkByUserIDLoginType(ctx, arg)
	m.queryLatencies.WithLabelValues("GetUserLinkByUserIDLoginType").Observe(time.Since(start).Seconds())
	return link, err
}

func (m metricsStore) GetUsers(ctx context.Context, arg database.GetUsersParams) ([]database.GetUsersRow, error) {
	start := time.Now()
	users, err := m.s.GetUsers(ctx, arg)
	m.queryLatencies.WithLabelValues("GetUsers").Observe(time.Since(start).Seconds())
	return users, err
}

func (m metricsStore) GetUsersByIDs(ctx context.Context, ids []uuid.UUID) ([]database.User, error) {
	start := time.Now()
	users, err := m.s.GetUsersByIDs(ctx, ids)
	m.queryLatencies.WithLabelValues("GetUsersByIDs").Observe(time.Since(start).Seconds())
	return users, err
}

func (m metricsStore) GetWorkspaceAgentByAuthToken(ctx context.Context, authToken uuid.UUID) (database.WorkspaceAgent, error) {
	start := time.Now()
	agent, err := m.s.GetWorkspaceAgentByAuthToken(ctx, authToken)
	m.queryLatencies.WithLabelValues("GetWorkspaceAgentByAuthToken").Observe(time.Since(start).Seconds())
	return agent, err
}

func (m metricsStore) GetWorkspaceAgentByID(ctx context.Context, id uuid.UUID) (database.WorkspaceAgent, error) {
	start := time.Now()
	agent, err := m.s.GetWorkspaceAgentByID(ctx, id)
	m.queryLatencies.WithLabelValues("GetWorkspaceAgentByID").Observe(time.Since(start).Seconds())
	return agent, err
}

func (m metricsStore) GetWorkspaceAgentByInstanceID(ctx context.Context, authInstanceID string) (database.WorkspaceAgent, error) {
	start := time.Now()
	agent, err := m.s.GetWorkspaceAgentByInstanceID(ctx, authInstanceID)
	m.queryLatencies.WithLabelValues("GetWorkspaceAgentByInstanceID").Observe(time.Since(start).Seconds())
	return agent, err
}

func (m metricsStore) GetWorkspaceAgentMetadata(ctx context.Context, workspaceAgentID uuid.UUID) ([]database.WorkspaceAgentMetadatum, error) {
	start := time.Now()
	metadata, err := m.s.GetWorkspaceAgentMetadata(ctx, workspaceAgentID)
	m.queryLatencies.WithLabelValues("GetWorkspaceAgentMetadata").Observe(time.Since(start).Seconds())
	return metadata, err
}

func (m metricsStore) GetWorkspaceAgentStartupLogsAfter(ctx context.Context, arg database.GetWorkspaceAgentStartupLogsAfterParams) ([]database.WorkspaceAgentStartupLog, error) {
	start := time.Now()
	logs, err := m.s.GetWorkspaceAgentStartupLogsAfter(ctx, arg)
	m.queryLatencies.WithLabelValues("GetWorkspaceAgentStartupLogsAfter").Observe(time.Since(start).Seconds())
	return logs, err
}

func (m metricsStore) GetWorkspaceAgentStats(ctx context.Context, createdAt time.Time) ([]database.GetWorkspaceAgentStatsRow, error) {
	start := time.Now()
	stats, err := m.s.GetWorkspaceAgentStats(ctx, createdAt)
	m.queryLatencies.WithLabelValues("GetWorkspaceAgentStats").Observe(time.Since(start).Seconds())
	return stats, err
}

func (m metricsStore) GetWorkspaceAgentStatsAndLabels(ctx context.Context, createdAt time.Time) ([]database.GetWorkspaceAgentStatsAndLabelsRow, error) {
	start := time.Now()
	stats, err := m.s.GetWorkspaceAgentStatsAndLabels(ctx, createdAt)
	m.queryLatencies.WithLabelValues("GetWorkspaceAgentStatsAndLabels").Observe(time.Since(start).Seconds())
	return stats, err
}

func (m metricsStore) GetWorkspaceAgentsByResourceIDs(ctx context.Context, ids []uuid.UUID) ([]database.WorkspaceAgent, error) {
	start := time.Now()
	agents, err := m.s.GetWorkspaceAgentsByResourceIDs(ctx, ids)
	m.queryLatencies.WithLabelValues("GetWorkspaceAgentsByResourceIDs").Observe(time.Since(start).Seconds())
	return agents, err
}

func (m metricsStore) GetWorkspaceAgentsCreatedAfter(ctx context.Context, createdAt time.Time) ([]database.WorkspaceAgent, error) {
	start := time.Now()
	agents, err := m.s.GetWorkspaceAgentsCreatedAfter(ctx, createdAt)
	m.queryLatencies.WithLabelValues("GetWorkspaceAgentsCreatedAfter").Observe(time.Since(start).Seconds())
	return agents, err
}

func (m metricsStore) GetWorkspaceAgentsInLatestBuildByWorkspaceID(ctx context.Context, workspaceID uuid.UUID) ([]database.WorkspaceAgent, error) {
	start := time.Now()
	agents, err := m.s.GetWorkspaceAgentsInLatestBuildByWorkspaceID(ctx, workspaceID)
	m.queryLatencies.WithLabelValues("GetWorkspaceAgentsInLatestBuildByWorkspaceID").Observe(time.Since(start).Seconds())
	return agents, err
}

func (m metricsStore) GetWorkspaceAppByAgentIDAndSlug(ctx context.Context, arg database.GetWorkspaceAppByAgentIDAndSlugParams) (database.WorkspaceApp, error) {
	start := time.Now()
	app, err := m.s.GetWorkspaceAppByAgentIDAndSlug(ctx, arg)
	m.queryLatencies.WithLabelValues("GetWorkspaceAppByAgentIDAndSlug").Observe(time.Since(start).Seconds())
	return app, err
}

func (m metricsStore) GetWorkspaceAppsByAgentID(ctx context.Context, agentID uuid.UUID) ([]database.WorkspaceApp, error) {
	start := time.Now()
	apps, err := m.s.GetWorkspaceAppsByAgentID(ctx, agentID)
	m.queryLatencies.WithLabelValues("GetWorkspaceAppsByAgentID").Observe(time.Since(start).Seconds())
	return apps, err
}

func (m metricsStore) GetWorkspaceAppsByAgentIDs(ctx context.Context, ids []uuid.UUID) ([]database.WorkspaceApp, error) {
	start := time.Now()
	apps, err := m.s.GetWorkspaceAppsByAgentIDs(ctx, ids)
	m.queryLatencies.WithLabelValues("GetWorkspaceAppsByAgentIDs").Observe(time.Since(start).Seconds())
	return apps, err
}

func (m metricsStore) GetWorkspaceAppsCreatedAfter(ctx context.Context, createdAt time.Time) ([]database.WorkspaceApp, error) {
	start := time.Now()
	apps, err := m.s.GetWorkspaceAppsCreatedAfter(ctx, createdAt)
	m.queryLatencies.WithLabelValues("GetWorkspaceAppsCreatedAfter").Observe(time.Since(start).Seconds())
	return apps, err
}

func (m metricsStore) GetWorkspaceBuildByID(ctx context.Context, id uuid.UUID) (database.WorkspaceBuild, error) {
	start := time.Now()
	build, err := m.s.GetWorkspaceBuildByID(ctx, id)
	m.queryLatencies.WithLabelValues("GetWorkspaceBuildByID").Observe(time.Since(start).Seconds())
	return build, err
}

func (m metricsStore) GetWorkspaceBuildByJobID(ctx context.Context, jobID uuid.UUID) (database.WorkspaceBuild, error) {
	start := time.Now()
	build, err := m.s.GetWorkspaceBuildByJobID(ctx, jobID)
	m.queryLatencies.WithLabelValues("GetWorkspaceBuildByJobID").Observe(time.Since(start).Seconds())
	return build, err
}

func (m metricsStore) GetWorkspaceBuildByWorkspaceIDAndBuildNumber(ctx context.Context, arg database.GetWorkspaceBuildByWorkspaceIDAndBuildNumberParams) (database.WorkspaceBuild, error) {
	start := time.Now()
	build, err := m.s.GetWorkspaceBuildByWorkspaceIDAndBuildNumber(ctx, arg)
	m.queryLatencies.WithLabelValues("GetWorkspaceBuildByWorkspaceIDAndBuildNumber").Observe(time.Since(start).Seconds())
	return build, err
}

func (m metricsStore) GetWorkspaceBuildParameters(ctx context.Context, workspaceBuildID uuid.UUID) ([]database.WorkspaceBuildParameter, error) {
	start := time.Now()
	params, err := m.s.GetWorkspaceBuildParameters(ctx, workspaceBuildID)
	m.queryLatencies.WithLabelValues("GetWorkspaceBuildParameters").Observe(time.Since(start).Seconds())
	return params, err
}

func (m metricsStore) GetWorkspaceBuildsByWorkspaceID(ctx context.Context, arg database.GetWorkspaceBuildsByWorkspaceIDParams) ([]database.WorkspaceBuild, error) {
	start := time.Now()
	builds, err := m.s.GetWorkspaceBuildsByWorkspaceID(ctx, arg)
	m.queryLatencies.WithLabelValues("GetWorkspaceBuildsByWorkspaceID").Observe(time.Since(start).Seconds())
	return builds, err
}

func (m metricsStore) GetWorkspaceBuildsCreatedAfter(ctx context.Context, createdAt time.Time) ([]database.WorkspaceBuild, error) {
	start := time.Now()
	builds, err := m.s.GetWorkspaceBuildsCreatedAfter(ctx, createdAt)
	m.queryLatencies.WithLabelValues("GetWorkspaceBuildsCreatedAfter").Observe(time.Since(start).Seconds())
	return builds, err
}

func (m metricsStore) GetWorkspaceByAgentID(ctx context.Context, agentID uuid.UUID) (database.Workspace, error) {
	start := time.Now()
	workspace, err := m.s.GetWorkspaceByAgentID(ctx, agentID)
	m.queryLatencies.WithLabelValues("GetWorkspaceByAgentID").Observe(time.Since(start).Seconds())
	return workspace, err
}

func (m metricsStore) GetWorkspaceByID(ctx context.Context, id uuid.UUID) (database.Workspace, error) {
	start := time.Now()
	workspace, err := m.s.GetWorkspaceByID(ctx, id)
	m.queryLatencies.WithLabelValues("GetWorkspaceByID").Observe(time.Since(start).Seconds())
	return workspace, err
}

func (m metricsStore) GetWorkspaceByOwnerIDAndName(ctx context.Context, arg database.GetWorkspaceByOwnerIDAndNameParams) (database.Workspace, error) {
	start := time.Now()
	workspace, err := m.s.GetWorkspaceByOwnerIDAndName(ctx, arg)
	m.queryLatencies.WithLabelValues("GetWorkspaceByOwnerIDAndName").Observe(time.Since(start).Seconds())
	return workspace, err
}

func (m metricsStore) GetWorkspaceByWorkspaceAppID(ctx context.Context, workspaceAppID uuid.UUID) (database.Workspace, error) {
	start := time.Now()
	workspace, err := m.s.GetWorkspaceByWorkspaceAppID(ctx, workspaceAppID)
	m.queryLatencies.WithLabelValues("GetWorkspaceByWorkspaceAppID").Observe(time.Since(start).Seconds())
	return workspace, err
}

func (m metricsStore) GetWorkspaceProxies(ctx context.Context) ([]database.WorkspaceProxy, error) {
	start := time.Now()
	proxies, err := m.s.GetWorkspaceProxies(ctx)
	m.queryLatencies.WithLabelValues("GetWorkspaceProxies").Observe(time.Since(start).Seconds())
	return proxies, err
}

func (m metricsStore) GetWorkspaceProxyByHostname(ctx context.Context, arg database.GetWorkspaceProxyByHostnameParams) (database.WorkspaceProxy, error) {
	start := time.Now()
	proxy, err := m.s.GetWorkspaceProxyByHostname(ctx, arg)
	m.queryLatencies.WithLabelValues("GetWorkspaceProxyByHostname").Observe(time.Since(start).Seconds())
	return proxy, err
}

func (m metricsStore) GetWorkspaceProxyByID(ctx context.Context, id uuid.UUID) (database.WorkspaceProxy, error) {
	start := time.Now()
	proxy, err := m.s.GetWorkspaceProxyByID(ctx, id)
	m.queryLatencies.WithLabelValues("GetWorkspaceProxyByID").Observe(time.Since(start).Seconds())
	return proxy, err
}

func (m metricsStore) GetWorkspaceProxyByName(ctx context.Context, name string) (database.WorkspaceProxy, error) {
	start := time.Now()
	proxy, err := m.s.GetWorkspaceProxyByName(ctx, name)
	m.queryLatencies.WithLabelValues("GetWorkspaceProxyByName").Observe(time.Since(start).Seconds())
	return proxy, err
}

func (m metricsStore) GetWorkspaceResourceByID(ctx context.Context, id uuid.UUID) (database.WorkspaceResource, error) {
	start := time.Now()
	resource, err := m.s.GetWorkspaceResourceByID(ctx, id)
	m.queryLatencies.WithLabelValues("GetWorkspaceResourceByID").Observe(time.Since(start).Seconds())
	return resource, err
}

func (m metricsStore) GetWorkspaceResourceMetadataByResourceIDs(ctx context.Context, ids []uuid.UUID) ([]database.WorkspaceResourceMetadatum, error) {
	start := time.Now()
	metadata, err := m.s.GetWorkspaceResourceMetadataByResourceIDs(ctx, ids)
	m.queryLatencies.WithLabelValues("GetWorkspaceResourceMetadataByResourceIDs").Observe(time.Since(start).Seconds())
	return metadata, err
}

func (m metricsStore) GetWorkspaceResourceMetadataCreatedAfter(ctx context.Context, createdAt time.Time) ([]database.WorkspaceResourceMetadatum, error) {
	start := time.Now()
	metadata, err := m.s.GetWorkspaceResourceMetadataCreatedAfter(ctx, createdAt)
	m.queryLatencies.WithLabelValues("GetWorkspaceResourceMetadataCreatedAfter").Observe(time.Since(start).Seconds())
	return metadata, err
}

func (m metricsStore) GetWorkspaceResourcesByJobID(ctx context.Context, jobID uuid.UUID) ([]database.WorkspaceResource, error) {
	start := time.Now()
	resources, err := m.s.GetWorkspaceResourcesByJobID(ctx, jobID)
	m.queryLatencies.WithLabelValues("GetWorkspaceResourcesByJobID").Observe(time.Since(start).Seconds())
	return resources, err
}

func (m metricsStore) GetWorkspaceResourcesByJobIDs(ctx context.Context, ids []uuid.UUID) ([]database.WorkspaceResource, error) {
	start := time.Now()
	resources, err := m.s.GetWorkspaceResourcesByJobIDs(ctx, ids)
	m.queryLatencies.WithLabelValues("GetWorkspaceResourcesByJobIDs").Observe(time.Since(start).Seconds())
	return resources, err
}

func (m metricsStore) GetWorkspaceResourcesCreatedAfter(ctx context.Context, createdAt time.Time) ([]database.WorkspaceResource, error) {
	start := time.Now()
	resources, err := m.s.GetWorkspaceResourcesCreatedAfter(ctx, createdAt)
	m.queryLatencies.WithLabelValues("GetWorkspaceResourcesCreatedAfter").Observe(time.Since(start).Seconds())
	return resources, err
}

func (m metricsStore) GetWorkspaces(ctx context.Context, arg database.GetWorkspacesParams) ([]database.GetWorkspacesRow, error) {
	start := time.Now()
	workspaces, err := m.s.GetWorkspaces(ctx, arg)
	m.queryLatencies.WithLabelValues("GetWorkspaces").Observe(time.Since(start).Seconds())
	return workspaces, err
}

func (m metricsStore) GetWorkspacesEligibleForAutoStartStop(ctx context.Context, now time.Time) ([]database.Workspace, error) {
	start := time.Now()
	workspaces, err := m.s.GetWorkspacesEligibleForAutoStartStop(ctx, now)
	m.queryLatencies.WithLabelValues("GetWorkspacesEligibleForAutoStartStop").Observe(time.Since(start).Seconds())
	return workspaces, err
}

func (m metricsStore) InsertAPIKey(ctx context.Context, arg database.InsertAPIKeyParams) (database.APIKey, error) {
	start := time.Now()
	key, err := m.s.InsertAPIKey(ctx, arg)
	m.queryLatencies.WithLabelValues("InsertAPIKey").Observe(time.Since(start).Seconds())
	return key, err
}

func (m metricsStore) InsertAllUsersGroup(ctx context.Context, organizationID uuid.UUID) (database.Group, error) {
	start := time.Now()
	group, err := m.s.InsertAllUsersGroup(ctx, organizationID)
	m.queryLatencies.WithLabelValues("InsertAllUsersGroup").Observe(time.Since(start).Seconds())
	return group, err
}

func (m metricsStore) InsertAuditLog(ctx context.Context, arg database.InsertAuditLogParams) (database.AuditLog, error) {
	start := time.Now()
	log, err := m.s.InsertAuditLog(ctx, arg)
	m.queryLatencies.WithLabelValues("InsertAuditLog").Observe(time.Since(start).Seconds())
	return log, err
}

func (m metricsStore) InsertDERPMeshKey(ctx context.Context, value string) error {
	start := time.Now()
	err := m.s.InsertDERPMeshKey(ctx, value)
	m.queryLatencies.WithLabelValues("InsertDERPMeshKey").Observe(time.Since(start).Seconds())
	return err
}

func (m metricsStore) InsertDeploymentID(ctx context.Context, value string) error {
	start := time.Now()
	err := m.s.InsertDeploymentID(ctx, value)
	m.queryLatencies.WithLabelValues("InsertDeploymentID").Observe(time.Since(start).Seconds())
	return err
}

func (m metricsStore) InsertFile(ctx context.Context, arg database.InsertFileParams) (database.File, error) {
	start := time.Now()
	file, err := m.s.InsertFile(ctx, arg)
	m.queryLatencies.WithLabelValues("InsertFile").Observe(time.Since(start).Seconds())
	return file, err
}

func (m metricsStore) InsertGitAuthLink(ctx context.Context, arg database.InsertGitAuthLinkParams) (database.GitAuthLink, error) {
	start := time.Now()
	link, err := m.s.InsertGitAuthLink(ctx, arg)
	m.queryLatencies.WithLabelValues("InsertGitAuthLink").Observe(time.Since(start).Seconds())
	return link, err
}

func (m metricsStore) InsertGitSSHKey(ctx context.Context, arg database.InsertGitSSHKeyParams) (database.GitSSHKey, error) {
	start := time.Now()
	key, err := m.s.InsertGitSSHKey(ctx, arg)
	m.queryLatencies.WithLabelValues("InsertGitSSHKey").Observe(time.Since(start).Seconds())
	return key, err
}

func (m metricsStore) InsertGroup(ctx context.Context, arg database.InsertGroupParams) (database.Group, error) {
	start := time.Now()
	group, err := m.s.InsertGroup(ctx, arg)
	m.queryLatencies.WithLabelValues("InsertGroup").Observe(time.Since(start).Seconds())
	return group, err
}

func (m metricsStore) InsertGroupMember(ctx context.Context, arg database.InsertGroupMemberParams) error {
	start := time.Now()
	err := m.s.InsertGroupMember(ctx, arg)
	m.queryLatencies.WithLabelValues("InsertGroupMember").Observe(time.Since(start).Seconds())
	return err
}

func (m metricsStore) InsertLicense(ctx context.Context, arg database.InsertLicenseParams) (database.License, error) {
	start := time.Now()
	license, err := m.s.InsertLicense(ctx, arg)
	m.queryLatencies.WithLabelValues("InsertLicense").Observe(time.Since(start).Seconds())
	return license, err
}

func (m metricsStore) InsertOrganization(ctx context.Context, arg database.InsertOrganizationParams) (database.Organization, error) {
	start := time.Now()
	organization, err := m.s.InsertOrganization(ctx, arg)
	m.queryLatencies.WithLabelValues("InsertOrganization").Observe(time.Since(start).Seconds())
	return organization, err
}

func (m metricsStore) InsertOrganizationMember(ctx context.Context, arg database.InsertOrganizationMemberParams) (database.OrganizationMember, error) {
	start := time.Now()
	member, err := m.s.InsertOrganizationMember(ctx, arg)
	m.queryLatencies.WithLabelValues("InsertOrganizationMember").Observe(time.Since(start).Seconds())
	return member, err
}

func (m metricsStore) InsertParameterSchema(ctx context.Context, arg database.InsertParameterSchemaParams) (database.ParameterSchema, error) {
	start := time.Now()
	schema, err := m.s.InsertParameterSchema(ctx, arg)
	m.queryLatencies.WithLabelValues("InsertParameterSchema").Observe(time.Since(start).Seconds())
	return schema, err
}

func (m metricsStore) InsertParameterValue(ctx context.Context, arg database.InsertParameterValueParams) (database.ParameterValue, error) {
	start := time.Now()
	value, err := m.s.InsertParameterValue(ctx, arg)
	m.queryLatencies.WithLabelValues("InsertParameterValue").Observe(time.Since(start).Seconds())
	return value, err
}

func (m metricsStore) InsertProvisionerDaemon(ctx context.Context, arg database.InsertProvisionerDaemonParams) (database.ProvisionerDaemon, error) {
	start := time.Now()
	daemon, err := m.s.InsertProvisionerDaemon(ctx, arg)
	m.queryLatencies.WithLabelValues("InsertProvisionerDaemon").Observe(time.Since(start).Seconds())
	return daemon, err
}

func (m metricsStore) InsertProvisionerJob(ctx context.Context, arg database.InsertProvisionerJobParams) (database.ProvisionerJob, error) {
	start := time.Now()
	job, err := m.s.InsertProvisionerJob(ctx, arg)
	m.queryLatencies.WithLabelValues("InsertProvisionerJob").Observe(time.Since(start).Seconds())
	return job, err
}

func (m metricsStore) InsertProvisionerJobLogs(ctx context.Context, arg database.InsertProvisionerJobLogsParams) ([]database.ProvisionerJobLog, error) {
	start := time.Now()
	logs, err := m.s.InsertProvisionerJobLogs(ctx, arg)
	m.queryLatencies.WithLabelValues("InsertProvisionerJobLogs").Observe(time.Since(start).Seconds())
	return logs, err
}

func (m metricsStore) InsertReplica(ctx context.Context, arg database.InsertReplicaParams) (database.Replica, error) {
	start := time.Now()
	replica, err := m.s.InsertReplica(ctx, arg)
	m.queryLatencies.WithLabelValues("InsertReplica").Observe(time.Since(start).Seconds())
	return replica, err
}

func (m metricsStore) InsertTemplate(ctx context.Context, arg database.InsertTemplateParams) (database.Template, error) {
	start := time.Now()
	template, err := m.s.InsertTemplate(ctx, arg)
	m.queryLatencies.WithLabelValues("InsertTemplate").Observe(time.Since(start).Seconds())
	return template, err
}

func (m metricsStore) InsertTemplateVersion(ctx context.Context, arg database.InsertTemplateVersionParams) (database.TemplateVersion, error) {
	start := time.Now()
	version, err := m.s.InsertTemplateVersion(ctx, arg)
	m.queryLatencies.WithLabelValues("InsertTemplateVersion").Observe(time.Since(start).Seconds())
	return version, err
}

func (m metricsStore) InsertTemplateVersionParameter(ctx context.Context, arg database.InsertTemplateVersionParameterParams) (database.TemplateVersionParameter, error) {
	start := time.Now()
	parameter, err := m.s.InsertTemplateVersionParameter(ctx, arg)
	m.queryLatencies.WithLabelValues("InsertTemplateVersionParameter").Observe(time.Since(start).Seconds())
	return parameter, err
}

func (m metricsStore) InsertTemplateVersionVariable(ctx context.Context, arg database.InsertTemplateVersionVariableParams) (database.TemplateVersionVariable, error) {
	start := time.Now()
	variable, err := m.s.InsertTemplateVersionVariable(ctx, arg)
	m.queryLatencies.WithLabelValues("InsertTemplateVersionVariable").Observe(time.Since(start).Seconds())
	return variable, err
}

func (m metricsStore) InsertUser(ctx context.Context, arg database.InsertUserParams) (database.User, error) {
	start := time.Now()
	user, err := m.s.InsertUser(ctx, arg)
	m.queryLatencies.WithLabelValues("InsertUser").Observe(time.Since(start).Seconds())
	return user, err
}

func (m metricsStore) InsertUserGroupsByName(ctx context.Context, arg database.InsertUserGroupsByNameParams) error {
	start := time.Now()
	err := m.s.InsertUserGroupsByName(ctx, arg)
	m.queryLatencies.WithLabelValues("InsertUserGroupsByName").Observe(time.Since(start).Seconds())
	return err
}

func (m metricsStore) InsertUserLink(ctx context.Context, arg database.InsertUserLinkParams) (database.UserLink, error) {
	start := time.Now()
	link, err := m.s.InsertUserLink(ctx, arg)
	m.queryLatencies.WithLabelValues("InsertUserLink").Observe(time.Since(start).Seconds())
	return link, err
}

func (m metricsStore) InsertWorkspace(ctx context.Context, arg database.InsertWorkspaceParams) (database.Workspace, error) {
	start := time.Now()
	workspace, err := m.s.InsertWorkspace(ctx, arg)
	m.queryLatencies.WithLabelValues("InsertWorkspace").Observe(time.Since(start).Seconds())
	return workspace, err
}

func (m metricsStore) InsertWorkspaceAgent(ctx context.Context, arg database.InsertWorkspaceAgentParams) (database.WorkspaceAgent, error) {
	start := time.Now()
	agent, err := m.s.InsertWorkspaceAgent(ctx, arg)
	m.queryLatencies.WithLabelValues("InsertWorkspaceAgent").Observe(time.Since(start).Seconds())
	return agent, err
}

func (m metricsStore) InsertWorkspaceAgentMetadata(ctx context.Context, arg database.InsertWorkspaceAgentMetadataParams) error {
	start := time.Now()
	err := m.s.InsertWorkspaceAgentMetadata(ctx, arg)
	m.queryLatencies.WithLabelValues("InsertWorkspaceAgentMetadata").Observe(time.Since(start).Seconds())
	return err
}

func (m metricsStore) InsertWorkspaceAgentStartupLogs(ctx context.Context, arg database.InsertWorkspaceAgentStartupLogsParams) ([]database.WorkspaceAgentStartupLog, error) {
	start := time.Now()
	logs, err := m.s.InsertWorkspaceAgentStartupLogs(ctx, arg)
	m.queryLatencies.WithLabelValues("InsertWorkspaceAgentStartupLogs").Observe(time.Since(start).Seconds())
	return logs, err
}

func (m metricsStore) InsertWorkspaceAgentStat(ctx context.Context, arg database.InsertWorkspaceAgentStatParams) (database.WorkspaceAgentStat, error) {
	start := time.Now()
	stat, err := m.s.InsertWorkspaceAgentStat(ctx, arg)
	m.queryLatencies.WithLabelValues("InsertWorkspaceAgentStat").Observe(time.Since(start).Seconds())
	return stat, err
}

func (m metricsStore) InsertWorkspaceApp(ctx context.Context, arg database.InsertWorkspaceAppParams) (database.WorkspaceApp, error) {
	start := time.Now()
	app, err := m.s.InsertWorkspaceApp(ctx, arg)
	m.queryLatencies.WithLabelValues("InsertWorkspaceApp").Observe(time.Since(start).Seconds())
	return app, err
}

func (m metricsStore) InsertWorkspaceBuild(ctx context.Context, arg database.InsertWorkspaceBuildParams) (database.WorkspaceBuild, error) {
	start := time.Now()
	build, err := m.s.InsertWorkspaceBuild(ctx, arg)
	m.queryLatencies.WithLabelValues("InsertWorkspaceBuild").Observe(time.Since(start).Seconds())
	return build, err
}

func (m metricsStore) InsertWorkspaceBuildParameters(ctx context.Context, arg database.InsertWorkspaceBuildParametersParams) error {
	start := time.Now()
	err := m.s.InsertWorkspaceBuildParameters(ctx, arg)
	m.queryLatencies.WithLabelValues("InsertWorkspaceBuildParameters").Observe(time.Since(start).Seconds())
	return err
}

func (m metricsStore) InsertWorkspaceProxy(ctx context.Context, arg database.InsertWorkspaceProxyParams) (database.WorkspaceProxy, error) {
	start := time.Now()
	proxy, err := m.s.InsertWorkspaceProxy(ctx, arg)
	m.queryLatencies.WithLabelValues("InsertWorkspaceProxy").Observe(time.Since(start).Seconds())
	return proxy, err
}

func (m metricsStore) InsertWorkspaceResource(ctx context.Context, arg database.InsertWorkspaceResourceParams) (database.WorkspaceResource, error) {
	start := time.Now()
	resource, err := m.s.InsertWorkspaceResource(ctx, arg)
	m.queryLatencies.WithLabelValues("InsertWorkspaceResource").Observe(time.Since(start).Seconds())
	return resource, err
}

func (m metricsStore) InsertWorkspaceResourceMetadata(ctx context.Context, arg database.InsertWorkspaceResourceMetadataParams) ([]database.WorkspaceResourceMetadatum, error) {
	start := time.Now()
	metadata, err := m.s.InsertWorkspaceResourceMetadata(ctx, arg)
	m.queryLatencies.WithLabelValues("InsertWorkspaceResourceMetadata").Observe(time.Since(start).Seconds())
	return metadata, err
}

func (m metricsStore) ParameterValue(ctx context.Context, id uuid.UUID) (database.ParameterValue, error) {
	start := time.Now()
	value, err := m.s.ParameterValue(ctx, id)
	m.queryLatencies.WithLabelValues("ParameterValue").Observe(time.Since(start).Seconds())
	return value, err
}

func (m metricsStore) ParameterValues(ctx context.Context, arg database.ParameterValuesParams) ([]database.ParameterValue, error) {
	start := time.Now()
	values, err := m.s.ParameterValues(ctx, arg)
	m.queryLatencies.WithLabelValues("ParameterValues").Observe(time.Since(start).Seconds())
	return values, err
}

func (m metricsStore) RegisterWorkspaceProxy(ctx context.Context, arg database.RegisterWorkspaceProxyParams) (database.WorkspaceProxy, error) {
	start := time.Now()
	proxy, err := m.s.RegisterWorkspaceProxy(ctx, arg)
	m.queryLatencies.WithLabelValues("RegisterWorkspaceProxy").Observe(time.Since(start).Seconds())
	return proxy, err
}

func (m metricsStore) TryAcquireLock(ctx context.Context, pgTryAdvisoryXactLock int64) (bool, error) {
	start := time.Now()
	ok, err := m.s.TryAcquireLock(ctx, pgTryAdvisoryXactLock)
	m.queryLatencies.WithLabelValues("TryAcquireLock").Observe(time.Since(start).Seconds())
	return ok, err
}

func (m metricsStore) UpdateAPIKeyByID(ctx context.Context, arg database.UpdateAPIKeyByIDParams) error {
	start := time.Now()
	err := m.s.UpdateAPIKeyByID(ctx, arg)
	m.queryLatencies.WithLabelValues("UpdateAPIKeyByID").Observe(time.Since(start).Seconds())
	return err
}

func (m metricsStore) UpdateGitAuthLink(ctx context.Context, arg database.UpdateGitAuthLinkParams) (database.GitAuthLink, error) {
	start := time.Now()
	link, err := m.s.UpdateGitAuthLink(ctx, arg)
	m.queryLatencies.WithLabelValues("UpdateGitAuthLink").Observe(time.Since(start).Seconds())
	return link, err
}

func (m metricsStore) UpdateGitSSHKey(ctx context.Context, arg database.UpdateGitSSHKeyParams) (database.GitSSHKey, error) {
	start := time.Now()
	key, err := m.s.UpdateGitSSHKey(ctx, arg)
	m.queryLatencies.WithLabelValues("UpdateGitSSHKey").Observe(time.Since(start).Seconds())
	return key, err
}

func (m metricsStore) UpdateGroupByID(ctx context.Context, arg database.UpdateGroupByIDParams) (database.Group, error) {
	start := time.Now()
	group, err := m.s.UpdateGroupByID(ctx, arg)
	m.queryLatencies.WithLabelValues("UpdateGroupByID").Observe(time.Since(start).Seconds())
	return group, err
}

func (m metricsStore) UpdateMemberRoles(ctx context.Context, arg database.UpdateMemberRolesParams) (database.OrganizationMember, error) {
	start := time.Now()
	member, err := m.s.UpdateMemberRoles(ctx, arg)
	m.queryLatencies.WithLabelValues("UpdateMemberRoles").Observe(time.Since(start).Seconds())
	return member, err
}

func (m metricsStore) UpdateProvisionerJobByID(ctx context.Context, arg database.UpdateProvisionerJobByIDParams) error {
	start := time.Now()
	err := m.s.UpdateProvisionerJobByID(ctx, arg)
	m.queryLatencies.WithLabelValues("UpdateProvisionerJobByID").Observe(time.Since(start).Seconds())
	return err
}

func (m metricsStore) UpdateProvisionerJobWithCancelByID(ctx context.Context, arg database.UpdateProvisionerJobWithCancelByIDParams) error {
	start := time.Now()
	err := m.s.UpdateProvisionerJobWithCancelByID(ctx, arg)
	m.queryLatencies.WithLabelValues("UpdateProvisionerJobWithCancelByID").Observe(time.Since(start).Seconds())
	return err
}

func (m metricsStore) UpdateProvisionerJobWithCompleteByID(ctx context.Context, arg database.UpdateProvisionerJobWithCompleteByIDParams) error {
	start := time.Now()
	err := m.s.UpdateProvisionerJobWithCompleteByID(ctx, arg)
	m.queryLatencies.WithLabelValues("UpdateProvisionerJobWithCompleteByID").Observe(time.Since(start).Seconds())
	return err
}

func (m metricsStore) UpdateReplica(ctx context.Context, arg database.UpdateReplicaParams) (database.Replica, error) {
	start := time.Now()
	replica, err := m.s.UpdateReplica(ctx, arg)
	m.queryLatencies.WithLabelValues("UpdateReplica").Observe(time.Since(start).Seconds())
	return replica, err
}

func (m metricsStore) UpdateTemplateACLByID(ctx context.Context, arg database.UpdateTemplateACLByIDParams) (database.Template, error) {
	start := time.Now()
	template, err := m.s.UpdateTemplateACLByID(ctx, arg)
	m.queryLatencies.WithLabelValues("UpdateTemplateACLByID").Observe(time.Since(start).Seconds())
	return template, err
}

func (m metricsStore) UpdateTemplateActiveVersionByID(ctx context.Context, arg database.UpdateTemplateActiveVersionByIDParams) error {
	start := time.Now()
	err := m.s.UpdateTemplateActiveVersionByID(ctx, arg)
	m.queryLatencies.WithLabelValues("UpdateTemplateActiveVersionByID").Observe(time.Since(start).Seconds())
	return err
}

func (m metricsStore) UpdateTemplateDeletedByID(ctx context.Context, arg database.UpdateTemplateDeletedByIDParams) error {
	start := time.Now()
	err := m.s.UpdateTemplateDeletedByID(ctx, arg)
	m.queryLatencies.WithLabelValues("UpdateTemplateDeletedByID").Observe(time.Since(start).Seconds())
	return err
}

func (m metricsStore) UpdateTemplateMetaByID(ctx context.Context, arg database.UpdateTemplateMetaByIDParams) (database.Template, error) {
	start := time.Now()
	template, err := m.s.UpdateTemplateMetaByID(ctx, arg)
	m.queryLatencies.WithLabelValues("UpdateTemplateMetaByID").Observe(time.Since(start).Seconds())
	return template, err
}

func (m metricsStore) UpdateTemplateScheduleByID(ctx context.Context, arg database.UpdateTemplateScheduleByIDParams) (database.Template, error) {
	start := time.Now()
	template, err := m.s.UpdateTemplateScheduleByID(ctx, arg)
	m.queryLatencies.WithLabelValues("UpdateTemplateScheduleByID").Observe(time.Since(start).Seconds())
	return template, err
}

func (m metricsStore) UpdateTemplateVersionByID(ctx context.Context, arg database.UpdateTemplateVersionByIDParams) (database.TemplateVersion, error) {
	start := time.Now()
	version, err := m.s.UpdateTemplateVersionByID(ctx, arg)
	m.queryLatencies.WithLabelValues("UpdateTemplateVersionByID").Observe(time.Since(start).Seconds())
	return version, err
}

func (m metricsStore) UpdateTemplateVersionDescriptionByJobID(ctx context.Context, arg database.UpdateTemplateVersionDescriptionByJobIDParams) error {
	start := time.Now()
	err := m.s.UpdateTemplateVersionDescriptionByJobID(ctx, arg)
	m.queryLatencies.WithLabelValues("UpdateTemplateVersionDescriptionByJobID").Observe(time.Since(start).Seconds())
	return err
}

func (m metricsStore) UpdateTemplateVersionGitAuthProvidersByJobID(ctx context.Context, arg database.UpdateTemplateVersionGitAuthProvidersByJobIDParams) error {
	start := time.Now()
	err := m.s.UpdateTemplateVersionGitAuthProvidersByJobID(ctx, arg)
	m.queryLatencies.WithLabelValues("UpdateTemplateVersionGitAuthProvidersByJobID").Observe(time.Since(start).Seconds())
	return err
}

func (m metricsStore) UpdateUserDeletedByID(ctx context.Context, arg database.UpdateUserDeletedByIDParams) error {
	start := time.Now()
	err := m.s.UpdateUserDeletedByID(ctx, arg)
	m.queryLatencies.WithLabelValues("UpdateUserDeletedByID").Observe(time.Since(start).Seconds())
	return err
}

func (m metricsStore) UpdateUserHashedPassword(ctx context.Context, arg database.UpdateUserHashedPasswordParams) error {
	start := time.Now()
	err := m.s.UpdateUserHashedPassword(ctx, arg)
	m.queryLatencies.WithLabelValues("UpdateUserHashedPassword").Observe(time.Since(start).Seconds())
	return err
}

func (m metricsStore) UpdateUserLastSeenAt(ctx context.Context, arg database.UpdateUserLastSeenAtParams) (database.User, error) {
	start := time.Now()
	user, err := m.s.UpdateUserLastSeenAt(ctx, arg)
	m.queryLatencies.WithLabelValues("UpdateUserLastSeenAt").Observe(time.Since(start).Seconds())
	return user, err
}

func (m metricsStore) UpdateUserLink(ctx context.Context, arg database.UpdateUserLinkParams) (database.UserLink, error) {
	start := time.Now()
	link, err := m.s.UpdateUserLink(ctx, arg)
	m.queryLatencies.WithLabelValues("UpdateUserLink").Observe(time.Since(start).Seconds())
	return link, err
}

func (m metricsStore) UpdateUserLinkedID(ctx context.Context, arg database.UpdateUserLinkedIDParams) (database.UserLink, error) {
	start := time.Now()
	link, err := m.s.UpdateUserLinkedID(ctx, arg)
	m.queryLatencies.WithLabelValues("UpdateUserLinkedID").Observe(time.Since(start).Seconds())
	return link, err
}

func (m metricsStore) UpdateUserProfile(ctx context.Context, arg database.UpdateUserProfileParams) (database.User, error) {
	start := time.Now()
	user, err := m.s.UpdateUserProfile(ctx, arg)
	m.queryLatencies.WithLabelValues("UpdateUserProfile").Observe(time.Since(start).Seconds())
	return user, err
}

func (m metricsStore) UpdateUserRoles(ctx context.Context, arg database.UpdateUserRolesParams) (database.User, error) {
	start := time.Now()
	user, err := m.s.UpdateUserRoles(ctx, arg)
	m.queryLatencies.WithLabelValues("UpdateUserRoles").Observe(time.Since(start).Seconds())
	return user, err
}

func (m metricsStore) UpdateUserStatus(ctx context.Context, arg database.UpdateUserStatusParams) (database.User, error) {
	start := time.Now()
	user, err := m.s.UpdateUserStatus(ctx, arg)
	m.queryLatencies.WithLabelValues("UpdateUserStatus").Observe(time.Since(start).Seconds())
	return user, err
}

func (m metricsStore) UpdateWorkspace(ctx context.Context, arg database.UpdateWorkspaceParams) (database.Workspace, error) {
	start := time.Now()
	workspace, err := m.s.UpdateWorkspace(ctx, arg)
	m.queryLatencies.WithLabelValues("UpdateWorkspace").Observe(time.Since(start).Seconds())
	return workspace, err
}

func (m metricsStore) UpdateWorkspaceAgentConnectionByID(ctx context.Context, arg database.UpdateWorkspaceAgentConnectionByIDParams) error {
	start := time.Now()
	err := m.s.UpdateWorkspaceAgentConnectionByID(ctx, arg)
	m.queryLatencies.WithLabelValues("UpdateWorkspaceAgentConnectionByID").Observe(time.Since(start).Seconds())
	return err
}

func (m metricsStore) UpdateWorkspaceAgentLifecycleStateByID(ctx context.Context, arg database.UpdateWorkspaceAgentLifecycleStateByIDParams) error {
	start := time.Now()
	err := m.s.UpdateWorkspaceAgentLifecycleStateByID(ctx, arg)
	m.queryLatencies.WithLabelValues("UpdateWorkspaceAgentLifecycleStateByID").Observe(time.Since(start).Seconds())
	return err
}

func (m metricsStore) UpdateWorkspaceAgentMetadata(ctx context.Context, arg database.UpdateWorkspaceAgentMetadataParams) error {
	start := time.Now()
	err := m.s.UpdateWorkspaceAgentMetadata(ctx, arg)
	m.queryLatencies.WithLabelValues("UpdateWorkspaceAgentMetadata").Observe(time.Since(start).Seconds())
	return err
}

func (m metricsStore) UpdateWorkspaceAgentStartupByID(ctx context.Context, arg database.UpdateWorkspaceAgentStartupByIDParams) error {
	start := time.Now()
	err := m.s.UpdateWorkspaceAgentStartupByID(ctx, arg)
	m.queryLatencies.WithLabelValues("UpdateWorkspaceAgentStartupByID").Observe(time.Since(start).Seconds())
	return err
}

func (m metricsStore) UpdateWorkspaceAgentStartupLogOverflowByID(ctx context.Context, arg database.UpdateWorkspaceAgentStartupLogOverflowByIDParams) error {
	start := time.Now()
	err := m.s.UpdateWorkspaceAgentStartupLogOverflowByID(ctx, arg)
	m.queryLatencies.WithLabelValues("UpdateWorkspaceAgentStartupLogOverflowByID").Observe(time.Since(start).Seconds())
	return err
}

func (m metricsStore) UpdateWorkspaceAppHealthByID(ctx context.Context, arg database.UpdateWorkspaceAppHealthByIDParams) error {
	start := time.Now()
	err := m.s.UpdateWorkspaceAppHealthByID(ctx, arg)
	m.queryLatencies.WithLabelValues("UpdateWorkspaceAppHealthByID").Observe(time.Since(start).Seconds())
	return err
}

func (m metricsStore) UpdateWorkspaceAutostart(ctx context.Context, arg database.UpdateWorkspaceAutostartParams) error {
	start := time.Now()
	err := m.s.UpdateWorkspaceAutostart(ctx, arg)
	m.queryLatencies.WithLabelValues("UpdateWorkspaceAutostart").Observe(time.Since(start).Seconds())
	return err
}

func (m metricsStore) UpdateWorkspaceBuildByID(ctx context.Context, arg database.UpdateWorkspaceBuildByIDParams) (database.WorkspaceBuild, error) {
	start := time.Now()
	build, err := m.s.UpdateWorkspaceBuildByID(ctx, arg)
	m.queryLatencies.WithLabelValues("UpdateWorkspaceBuildByID").Observe(time.Since(start).Seconds())
	return build, err
}

func (m metricsStore) UpdateWorkspaceBuildCostByID(ctx context.Context, arg database.UpdateWorkspaceBuildCostByIDParams) (database.WorkspaceBuild, error) {
	start := time.Now()
	build, err := m.s.UpdateWorkspaceBuildCostByID(ctx, arg)
	m.queryLatencies.WithLabelValues("UpdateWorkspaceBuildCostByID").Observe(time.Since(start).Seconds())
	return build, err
}

func (m metricsStore) UpdateWorkspaceDeletedByID(ctx context.Context, arg database.UpdateWorkspaceDeletedByIDParams) error {
	start := time.Now()
	err := m.s.UpdateWorkspaceDeletedByID(ctx, arg)
	m.queryLatencies.WithLabelValues("UpdateWorkspaceDeletedByID").Observe(time.Since(start).Seconds())
	return err
}

func (m metricsStore) UpdateWorkspaceLastUsedAt(ctx context.Context, arg database.UpdateWorkspaceLastUsedAtParams) error {
	start := time.Now()
	err := m.s.UpdateWorkspaceLastUsedAt(ctx, arg)
	m.queryLatencies.WithLabelValues("UpdateWorkspaceLastUsedAt").Observe(time.Since(start).Seconds())
	return err
}

func (m metricsStore) UpdateWorkspaceProxy(ctx context.Context, arg database.UpdateWorkspaceProxyParams) (database.WorkspaceProxy, error) {
	start := time.Now()
	proxy, err := m.s.UpdateWorkspaceProxy(ctx, arg)
	m.queryLatencies.WithLabelValues("UpdateWorkspaceProxy").Observe(time.Since(start).Seconds())
	return proxy, err
}

func (m metricsStore) UpdateWorkspaceProxyDeleted(ctx context.Context, arg database.UpdateWorkspaceProxyDeletedParams) error {
	start := time.Now()
	err := m.s.UpdateWorkspaceProxyDeleted(ctx, arg)
	m.queryLatencies.WithLabelValues("UpdateWorkspaceProxyDeleted").Observe(time.Since(start).Seconds())
	return err
}

func (m metricsStore) UpdateWorkspaceTTL(ctx context.Context, arg database.UpdateWorkspaceTTLParams) error {
	start := time.Now()
	err := m.s.UpdateWorkspaceTTL(ctx, arg)
	m.queryLatencies.WithLabelValues("UpdateWorkspaceTTL").Observe(time.Since(start).Seconds())
	return err
}

func (m metricsStore) UpdateWorkspaceTTLToBeWithinTemplateMax(ctx context.Context, arg database.UpdateWorkspaceTTLToBeWithinTemplateMaxParams) error {
	start := time.Now()
	err := m.s.UpdateWorkspaceTTLToBeWithinTemplateMax(ctx, arg)
	m.queryLatencies.WithLabelValues("UpdateWorkspaceTTLToBeWithinTemplateMax").Observe(time.Since(start).Seconds())
	return err
}

func (m metricsStore) UpsertAppSecurityKey(ctx context.Context, value string) error {
	start := time.Now()
	err := m.s.UpsertAppSecurityKey(ctx, value)
	m.queryLatencies.WithLabelValues("UpsertAppSecurityKey").Observe(time.Since(start).Seconds())
	return err
}

func (m metricsStore) UpsertLastUpdateCheck(ctx context.Context, value string) error {
	start := time.Now()
	err := m.s.UpsertLastUpdateCheck(ctx, value)
	m.queryLatencies.WithLabelValues("UpsertLastUpdateCheck").Observe(time.Since(start).Seconds())
	return err
}

func (m metricsStore) UpsertLogoURL(ctx context.Context, value string) error {
	start := time.Now()
	err := m.s.UpsertLogoURL(ctx, value)
	m.queryLatencies.WithLabelValues("UpsertLogoURL").Observe(time.Since(start).Seconds())
	return err
}

func (m metricsStore) UpsertServiceBanner(ctx context.Context, value string) error {
	start := time.Now()
	err := m.s.UpsertServiceBanner(ctx, value)
	m.queryLatencies.WithLabelValues("UpsertServiceBanner").Observe(time.Since(start).Seconds())
	return err
}

func (m metricsStore) GetAuthorizedTemplates(ctx context.Context, arg database.GetTemplatesWithFilterParams, prepared rbac.PreparedAuthorized) ([]database.Template, error) {
	start := time.Now()
	templates, err := m.s.GetAuthorizedTemplates(ctx, arg, prepared)
	m.queryLatencies.WithLabelValues("GetAuthorizedTemplates").Observe(time.Since(start).Seconds())
	return templates, err
}

func (m metricsStore) GetTemplateGroupRoles(ctx context.Context, id uuid.UUID) ([]database.TemplateGroup, error) {
	start := time.Now()
	roles, err := m.s.GetTemplateGroupRoles(ctx, id)
	m.queryLatencies.WithLabelValues("GetTemplateGroupRoles").Observe(time.Since(start).Seconds())
	return roles, err
}

func (m metricsStore) GetTemplateUserRoles(ctx context.Context, id uuid.UUID) ([]database.TemplateUser, error) {
	start := time.Now()
	roles, err := m.s.GetTemplateUserRoles(ctx, id)
	m.queryLatencies.WithLabelValues("GetTemplateUserRoles").Observe(time.Since(start).Seconds())
	return roles, err
}

func (m metricsStore) GetAuthorizedWorkspaces(ctx context.Context, arg database.GetWorkspacesParams, prepared rbac.PreparedAuthorized) ([]database.GetWorkspacesRow, error) {
	start := time.Now()
	workspaces, err := m.s.GetAuthorizedWorkspaces(ctx, arg, prepared)
	m.queryLatencies.WithLabelValues("GetAuthorizedWorkspaces").Observe(time.Since(start).Seconds())
	return workspaces, err
}

func (m metricsStore) GetAuthorizedUserCount(ctx context.Context, arg database.GetFilteredUserCountParams, prepared rbac.PreparedAuthorized) (int64, error) {
	start := time.Now()
	count, err := m.s.GetAuthorizedUserCount(ctx, arg, prepared)
	m.queryLatencies.WithLabelValues("GetAuthorizedUserCount").Observe(time.Since(start).Seconds())
	return count, err
}