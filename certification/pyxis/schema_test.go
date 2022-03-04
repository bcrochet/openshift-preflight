package pyxis

var Schema = `
type Query {
  # Get container image by ID.
  get_image(id: String): ContainerImageResponse

  # Get container images. Exclude total for improved performance.
  find_images(
    sort_by: [SortBy]
    page: Int = 0
    page_size: Int = 50
    filter: ContainerImageFilter
  ): ContainerImagePaginatedResponse
  # Get a certification project test result
  get_cert_project_test_result(id: String): CertProjectTestResultResponse

  # Get many certification project test results
  get_cert_project_test_results(
    id: String
    sort_by: [SortBy]
    page: Int = 0
    page_size: Int = 50
    filter: CertProjectTestResultFilter
  ): CertProjectTestResultPaginatedResponse

  # Get test results by container image id
  get_test_results_by_image_id(
    id: String
    sort_by: [SortBy]
    page: Int = 0
    page_size: Int = 50
    filter: CertProjectTestResultFilter
  ): CertProjectTestResultPaginatedResponse

  # Get a certification project artifact
  get_cert_project_artifact(id: String): CertProjectArtifactResponse

  # Get a certification project artifacts
  get_cert_project_artifacts(
    id: String
    sort_by: [SortBy]
    page: Int = 0
    page_size: Int = 50
    filter: CertProjectArtifactFilter
  ): CertProjectArtifactPaginatedResponse

  # Get artifacts by container image id
  get_certification_artifacts_by_image_id(
    id: String
    sort_by: [SortBy]
    page: Int = 0
    page_size: Int = 50
    filter: CertProjectArtifactFilter
  ): CertProjectArtifactPaginatedResponse

  # Get an RPM manifest by ID
  get_rpm_manifest(id: String): ContainerImageRPMManifestResponse

  # Get the RPM manifest for an image
  get_image_rpm_manifest(id: String): ContainerImageRPMManifestResponse

  # Get certification project using its ID.
  get_certification_project(id: String): CertificationProjectResponse

  # Get certification project using Red Hat Connect project ID.
  get_certification_project_by_pid(pid: String): CertificationProjectResponse

  # Get images for certification project using its ID.
  find_certification_projects_images(
    id: String
    sort_by: [SortBy]
    page: Int = 0
    page_size: Int = 50
    filter: ContainerImageFilter
  ): ContainerImagePaginatedResponse
}

type StringResponse {
  data: String
  error: ResponseError
}

type ResponseError {
  status: Int
  detail: String
}

type ForwarderStatusResponse {
  data: ForwarderStatus
  error: ResponseError
}

#
type ForwarderStatus {
  forwarders: Forwarders
}

# Object with all log forwarder statuses
type Forwarders {
  fluentd: ForwarderStatusInfo
}

# Splunk forwarder status information
type ForwarderStatusInfo {
  # Forwarder status
  status: Boolean
}

type ContainerImageResponse {
  data: ContainerImage
  error: ResponseError
}

# Metadata about images contained in RedHat and ISV repositories
type ContainerImage {
  # The field contains an architecture for which the container image was built for. Value is used to distinguish between the default x86-64 architecture and other architectures. If the value is not set, the image was built for the x86-64 architecture.
  architecture: String

  # Brew related metadata.
  brew: Brew
  certifications: [Certification]
    @deprecated(
      reason: "The field is no longer supported. Certification test results were moved to test-results endpoint."
    )

  # A list of all content sets (YUM repositories) from where an image RPM content is.
  content_sets: [String]

  # A mapping of applicable advisories to RPM NEVRA. This data is required for scoring.
  cpe_ids: [String]

  # A mapping of applicable advisories for the base_images from the Red Hat repositories.
  cpe_ids_rh_base_images: [String]

  # Docker Image Digest. For Docker 1.10+ this is also known as the 'manifest digest'.
  docker_image_digest: String

  # Docker Image ID. For Docker 1.10+ this is also known as the 'config digest'.
  docker_image_id: String

  # The grade based on applicable updates and time provided by PST CVE engine.
  freshness_grades: [FreshnessGrade]
  object_type: String

  # Data parsed from image metadata.
  # These fields are not computed from any other source.
  parsed_data: ParsedData

  # Published repositories associated with the container image.
  repositories: [ContainerImageRepo]

  # The certification scan status. The field is generated based on certification info and it can't be used in query filter.
  scan_status: String @deprecated(reason: "The field is no longer supported.")

  # Indication if the image was certified.
  certified: Boolean

  # Indicates that an image was removed. Only unpublished images can be removed.
  deleted: Boolean

  # Image manifest digest.
  # Be careful, as this value is not unique among container image entries, as one image can be references several times.
  image_id: String

  # ID of the project in for ISV repositories. The ID can be also used to connect vendor to the image.
  isv_pid: String

  # The total size of the sum of all layers for each image in bytes. This is computed externally and may not match what is reported by the image metadata (see parsed_data.size).
  sum_layer_size_bytes: Int

  # Field for multiarch primary key
  top_layer_id: String

  # Hash (sha256) of the uncompressed top layer for this image (should be same value as - parsed_data.uncompressed_layer_sizes.0.layer_id)
  uncompressed_top_layer_id: String

  # Raw image configuration, such as output from docker inspect.
  raw_config: String

  # MongoDB unique _id
  _id: String

  # The date when the entry was created. Value is created automatically on creation.
  creation_date: DateTime

  # The date when the entry was last updated.
  last_update_date: DateTime
  edges: ContainerImageEdges
}

# Brew Build System related metadata.
type Brew {
  # Unique and immutable Brew build ID.
  build: String

  # Timestamp from Brew when the image has been succesfully built.
  completion_date: DateTime

  # Multi-Arch primary key.
  nvra: String

  # A package name in Brew.
  package: String
}

scalar DateTime

#
type Certification {
  assessment: [Assessment]
}

#
type Assessment {
  # Assesment name.
  name: String

  # Indicates if the assessment is required for certification.
  required_for_certification: Boolean

  # Indicates if the assesment was passed, True means yes.
  value: Boolean
}

# Grade based on applicable updates and time provided by PST CVE engine.
type FreshnessGrade {
  # Date after which the grade is no longer valid. See start_date for when the grade was effective. If no value is set, the grade applies forever. This should happen only for a grade of A (no vulnerabilities) or grade F.
  end_date: DateTime

  # The grade.
  grade: String

  # Date when the grade was added by the vulnerability engine.
  creation_date: DateTime

  # Date from which the grade is in effect. The grade is effective until the end_date, if end_date is set.
  start_date: DateTime
}

#
type ParsedData {
  architecture: String
  author: String
  command: String
  comment: String
  container: String

  # The 'created' date reported by image metadata. Stored as String because we do not have control on that format.
  created: String
  docker_image_digest: String
  docker_image_id: String

  # Version of docker reported by 'docker inspect' for this image.
  docker_version: String
  env_variables: [String]
  image_id: String
  labels: [Label]

  # Layer digests from the image.
  layers: [String]
  os: String
  ports: String

  # Repositories defined within an image as reported by yum command.
  repos: [ParsedDataRepo]

  # Size of this image as reported by image metadata.
  size: Int

  # Information about uncompressed layer sizes.
  uncompressed_layer_sizes: [UncompressedLayerSize]

  # Uncompressed images size in bytes (sum of uncompressed layers size).
  uncompressed_size_bytes: Int

  # The user on the images.
  user: String

  # Virtual size of this image as reported by image metadata.
  virtual_size: Int
}

# Image label.
type Label {
  # The name of the label
  name: String

  # Value of the label.
  value: String
}

#
type ParsedDataRepo {
  baseurl: String
  expire: String
  filename: String
  id: String
  name: String
  pkgs: String
  size: String
  updated: String
}

#
type UncompressedLayerSize {
  # The SHA256 layer ID.
  layer_id: String

  # The uncompressed layer size in bytes.
  size_bytes: Int
}

#
type ContainerImageRepo {
  # Store information about image comparison.
  comparison: ContainerImageRepoComparison

  # The _id's of the redHatContainerAdvisory that contains the content advisories.
  content_advisory_ids: [String]

  # The _id of the redHatContainerAdvisory that contains the image advisory.
  image_advisory_id: String

  # Available for multiarch images.
  manifest_list_digest: String

  # Available for single arch images.
  manifest_schema2_digest: String

  # Indicate if the image has been published to the container catalog.
  published: Boolean

  # Date the image was published to the container catalog.
  published_date: DateTime

  # When the image was pushed to this repository. For RH images this is picked from first found of advisory ship_date, brew completion_date, and finally repositories publish_date. For ISV images this TBD but is probably going to be only sourced from publish_date but could come from parsed_data.created.
  push_date: DateTime

  # Hostname of the registry where the repository can be accessed.
  registry: String

  # Repository name.
  repository: String

  # Image signing info.
  signatures: [SignatureInfo]

  # List of container tags assigned to this layer.
  tags: [ContainerImageRepoTag]
  edges: ContainerImageRepoEdges
}

#
type ContainerImageRepoComparison {
  # Mapping of a NVRA to multiple advisories IDs.
  advisory_rpm_mapping: [ContainerImageRepoComparisonMapping]

  # Reason why 'with_nvr' is or is not null.
  reason: String

  # Human readable reason.
  reason_text: String

  # List of rpms grouped by category (new, remove, upgrade, downgrade).
  rpms: ContainerImageRepoComparisonRPMs

  # NVR of image which this image was compared with.
  with_nvr: String
}

#
type ContainerImageRepoComparisonMapping {
  # Content advisory ID.
  advisory_ids: [String]

  # NVRA of the RPM related to advisories.
  nvra: String
}

#
type ContainerImageRepoComparisonRPMs {
  # List of NVRA which were downgraded in this image.
  downgrade: [String]

  # List of NVRA which were added to this image.
  new: [String]

  # List of NVRA which were removed in this image.
  remove: [String]

  # List of NVRA which were upgraded in this image.
  upgrade: [String]
}

#
type SignatureInfo {
  # The long 16-byte gpg key id.
  key_long_id: String

  # List of image tags that are signed with the given key.
  tags: [String]
}

#
type ContainerImageRepoTag {
  added_date: DateTime

  # Available when manifest_schema2_digest is not. All legacy images.
  manifest_schema1_digest: String

  # The name of the tag.
  name: String

  # Date this tag was removed from the image in this repo. If the tag is added back, add a new entry in 'tags' array.
  removed_date: DateTime
  edges: ContainerImageRepoTagEdges
}

type ContainerImageRepoTagEdges {
  tag_history: ContainerTagHistoryResponse
}

type ContainerTagHistoryResponse {
  data: ContainerTagHistory
  error: ResponseError
}

# The tag history stores a list of image that still have or used to have the given tag
type ContainerTagHistory {
  object_type: String

  # Hostname of the registry where the repository can be accessed.
  registry: String

  # Repository name.
  repository: String

  # The image tag name.
  tag: String

  # The tag type i.e. for floating or persistent.
  tag_type: String

  # Array with the tag history information.
  history: [History]

  # MongoDB unique _id
  _id: String

  # The date when the entry was created. Value is created automatically on creation.
  creation_date: DateTime

  # The date when the entry was last updated.
  last_update_date: DateTime
  edges: ContainerTagHistoryEdges
}

# Array with the tag history information.
type History {
  # Unique immutable build identifier in the brew build system.
  brew_build: String
  end_date: DateTime

  # The date for when the tag for the given docker_image_digest starts.
  start_date: DateTime
}

type ContainerTagHistoryEdges {
  images(
    page_size: Int = 50
    page: Int = 0
    filter: ContainerImageFilter
    sort_by: [SortBy]
  ): ContainerImagePaginatedResponse
}

type ContainerImagePaginatedResponse {
  data: [ContainerImage]
  error: ResponseError
  page: Int
  page_size: Int
  total: Int
}

input ContainerImageFilter {
  architecture: StringFilter
  brew: BrewFilter
  certifications_size: ListSizeFilter
  certifications_index: StringIndexFilter
  certifications_elemMatch: CertificationElemMatchFilter
  certifications: CertificationFilter
  content_sets_size: ListSizeFilter
  content_sets_index: StringIndexFilter
  content_sets: StringFilter
  cpe_ids_size: ListSizeFilter
  cpe_ids_index: StringIndexFilter
  cpe_ids: StringFilter
  cpe_ids_rh_base_images_size: ListSizeFilter
  cpe_ids_rh_base_images_index: StringIndexFilter
  cpe_ids_rh_base_images: StringFilter
  docker_image_digest: StringFilter
  docker_image_id: StringFilter
  freshness_grades_size: ListSizeFilter
  freshness_grades_index: StringIndexFilter
  freshness_grades_elemMatch: FreshnessGradeElemMatchFilter
  freshness_grades: FreshnessGradeFilter
  object_type: StringFilter
  parsed_data: ParsedDataFilter
  repositories_size: ListSizeFilter
  repositories_index: StringIndexFilter
  repositories_elemMatch: ContainerImageRepoElemMatchFilter
  repositories: ContainerImageRepoFilter
  scan_status: StringFilter
  certified: BooleanFilter
  deleted: BooleanFilter
  image_id: StringFilter
  isv_pid: StringFilter
  sum_layer_size_bytes: IntFilter
  top_layer_id: StringFilter
  uncompressed_top_layer_id: StringFilter
  raw_config: StringFilter
  _id: StringFilter
  creation_date: DateTimeFilter
  last_update_date: DateTimeFilter
  and: [ContainerImageFilter]
  or: [ContainerImageFilter]
  nor: [ContainerImageFilter]
  not: ContainerImageFilter
}

input StringFilter {
  eq: String
  ne: String
  le: String
  lt: String
  ge: String
  gt: String
  in: [String]
  out: [String]
  all: [String]
  iregex: String
  regex: String
  size: Int
}

input BrewFilter {
  build: StringFilter
  completion_date: DateTimeFilter
  nvra: StringFilter
  package: StringFilter
  and: [BrewFilter]
  or: [BrewFilter]
  nor: [BrewFilter]
  not: BrewFilter
}

input DateTimeFilter {
  eq: DateTime
  ne: DateTime
  le: DateTime
  lt: DateTime
  ge: DateTime
  gt: DateTime
  in: [DateTime]
}

input ListSizeFilter {
  eq: Int
}

input StringIndexFilter {
  condition: StringFilter
  index: Int
}

input CertificationElemMatchFilter {
  and: [CertificationFilter]
  or: [CertificationFilter]
  nor: [CertificationFilter]
  not: CertificationFilter
}

input CertificationFilter {
  assessment_size: ListSizeFilter
  assessment_index: StringIndexFilter
  assessment_elemMatch: AssessmentElemMatchFilter
  assessment: AssessmentFilter
  and: [CertificationFilter]
  or: [CertificationFilter]
  nor: [CertificationFilter]
  not: CertificationFilter
}

input AssessmentElemMatchFilter {
  and: [AssessmentFilter]
  or: [AssessmentFilter]
  nor: [AssessmentFilter]
  not: AssessmentFilter
}

input AssessmentFilter {
  name: StringFilter
  required_for_certification: BooleanFilter
  value: BooleanFilter
  and: [AssessmentFilter]
  or: [AssessmentFilter]
  nor: [AssessmentFilter]
  not: AssessmentFilter
}

input BooleanFilter {
  eq: Boolean
  ne: Boolean
}

input FreshnessGradeElemMatchFilter {
  and: [FreshnessGradeFilter]
  or: [FreshnessGradeFilter]
  nor: [FreshnessGradeFilter]
  not: FreshnessGradeFilter
}

input FreshnessGradeFilter {
  end_date: DateTimeFilter
  grade: StringFilter
  creation_date: DateTimeFilter
  start_date: DateTimeFilter
  and: [FreshnessGradeFilter]
  or: [FreshnessGradeFilter]
  nor: [FreshnessGradeFilter]
  not: FreshnessGradeFilter
}

input ParsedDataFilter {
  architecture: StringFilter
  author: StringFilter
  command: StringFilter
  comment: StringFilter
  container: StringFilter
  created: StringFilter
  docker_image_digest: StringFilter
  docker_image_id: StringFilter
  docker_version: StringFilter
  env_variables_size: ListSizeFilter
  env_variables_index: StringIndexFilter
  env_variables: StringFilter
  image_id: StringFilter
  labels_size: ListSizeFilter
  labels_index: StringIndexFilter
  labels_elemMatch: LabelElemMatchFilter
  labels: LabelFilter
  layers_size: ListSizeFilter
  layers_index: StringIndexFilter
  layers: StringFilter
  os: StringFilter
  ports: StringFilter
  repos_size: ListSizeFilter
  repos_index: StringIndexFilter
  repos_elemMatch: ParsedDataRepoElemMatchFilter
  repos: ParsedDataRepoFilter
  size: IntFilter
  uncompressed_layer_sizes_size: ListSizeFilter
  uncompressed_layer_sizes_index: StringIndexFilter
  uncompressed_layer_sizes_elemMatch: UncompressedLayerSizeElemMatchFilter
  uncompressed_layer_sizes: UncompressedLayerSizeFilter
  uncompressed_size_bytes: IntFilter
  user: StringFilter
  virtual_size: IntFilter
  and: [ParsedDataFilter]
  or: [ParsedDataFilter]
  nor: [ParsedDataFilter]
  not: ParsedDataFilter
}

input LabelElemMatchFilter {
  and: [LabelFilter]
  or: [LabelFilter]
  nor: [LabelFilter]
  not: LabelFilter
}

input LabelFilter {
  name: StringFilter
  value: StringFilter
  and: [LabelFilter]
  or: [LabelFilter]
  nor: [LabelFilter]
  not: LabelFilter
}

input ParsedDataRepoElemMatchFilter {
  and: [ParsedDataRepoFilter]
  or: [ParsedDataRepoFilter]
  nor: [ParsedDataRepoFilter]
  not: ParsedDataRepoFilter
}

input ParsedDataRepoFilter {
  baseurl: StringFilter
  expire: StringFilter
  filename: StringFilter
  id: StringFilter
  name: StringFilter
  pkgs: StringFilter
  size: StringFilter
  updated: StringFilter
  and: [ParsedDataRepoFilter]
  or: [ParsedDataRepoFilter]
  nor: [ParsedDataRepoFilter]
  not: ParsedDataRepoFilter
}

input IntFilter {
  eq: Int
  ne: Int
  le: Int
  lt: Int
  ge: Int
  gt: Int
  in: [Int]
}

input UncompressedLayerSizeElemMatchFilter {
  and: [UncompressedLayerSizeFilter]
  or: [UncompressedLayerSizeFilter]
  nor: [UncompressedLayerSizeFilter]
  not: UncompressedLayerSizeFilter
}

input UncompressedLayerSizeFilter {
  layer_id: StringFilter
  size_bytes: IntFilter
  and: [UncompressedLayerSizeFilter]
  or: [UncompressedLayerSizeFilter]
  nor: [UncompressedLayerSizeFilter]
  not: UncompressedLayerSizeFilter
}

input ContainerImageRepoElemMatchFilter {
  and: [ContainerImageRepoFilter]
  or: [ContainerImageRepoFilter]
  nor: [ContainerImageRepoFilter]
  not: ContainerImageRepoFilter
}

input ContainerImageRepoFilter {
  comparison: ContainerImageRepoComparisonFilter
  content_advisory_ids_size: ListSizeFilter
  content_advisory_ids_index: StringIndexFilter
  content_advisory_ids: StringFilter
  image_advisory_id: StringFilter
  manifest_list_digest: StringFilter
  manifest_schema2_digest: StringFilter
  published: BooleanFilter
  published_date: DateTimeFilter
  push_date: DateTimeFilter
  registry: StringFilter
  repository: StringFilter
  signatures_size: ListSizeFilter
  signatures_index: StringIndexFilter
  signatures_elemMatch: SignatureInfoElemMatchFilter
  signatures: SignatureInfoFilter
  tags_size: ListSizeFilter
  tags_index: StringIndexFilter
  tags_elemMatch: ContainerImageRepoTagElemMatchFilter
  tags: ContainerImageRepoTagFilter
  and: [ContainerImageRepoFilter]
  or: [ContainerImageRepoFilter]
  nor: [ContainerImageRepoFilter]
  not: ContainerImageRepoFilter
}

input ContainerImageRepoComparisonFilter {
  advisory_rpm_mapping_size: ListSizeFilter
  advisory_rpm_mapping_index: StringIndexFilter
  advisory_rpm_mapping_elemMatch: ContainerImageRepoComparisonMappingElemMatchFilter
  advisory_rpm_mapping: ContainerImageRepoComparisonMappingFilter
  reason: StringFilter
  reason_text: StringFilter
  rpms: ContainerImageRepoComparisonRPMsFilter
  with_nvr: StringFilter
  and: [ContainerImageRepoComparisonFilter]
  or: [ContainerImageRepoComparisonFilter]
  nor: [ContainerImageRepoComparisonFilter]
  not: ContainerImageRepoComparisonFilter
}

input ContainerImageRepoComparisonMappingElemMatchFilter {
  and: [ContainerImageRepoComparisonMappingFilter]
  or: [ContainerImageRepoComparisonMappingFilter]
  nor: [ContainerImageRepoComparisonMappingFilter]
  not: ContainerImageRepoComparisonMappingFilter
}

input ContainerImageRepoComparisonMappingFilter {
  advisory_ids_size: ListSizeFilter
  advisory_ids_index: StringIndexFilter
  advisory_ids: StringFilter
  nvra: StringFilter
  and: [ContainerImageRepoComparisonMappingFilter]
  or: [ContainerImageRepoComparisonMappingFilter]
  nor: [ContainerImageRepoComparisonMappingFilter]
  not: ContainerImageRepoComparisonMappingFilter
}

input ContainerImageRepoComparisonRPMsFilter {
  downgrade_size: ListSizeFilter
  downgrade_index: StringIndexFilter
  downgrade: StringFilter
  new_size: ListSizeFilter
  new_index: StringIndexFilter
  new: StringFilter
  remove_size: ListSizeFilter
  remove_index: StringIndexFilter
  remove: StringFilter
  upgrade_size: ListSizeFilter
  upgrade_index: StringIndexFilter
  upgrade: StringFilter
  and: [ContainerImageRepoComparisonRPMsFilter]
  or: [ContainerImageRepoComparisonRPMsFilter]
  nor: [ContainerImageRepoComparisonRPMsFilter]
  not: ContainerImageRepoComparisonRPMsFilter
}

input SignatureInfoElemMatchFilter {
  and: [SignatureInfoFilter]
  or: [SignatureInfoFilter]
  nor: [SignatureInfoFilter]
  not: SignatureInfoFilter
}

input SignatureInfoFilter {
  key_long_id: StringFilter
  tags_size: ListSizeFilter
  tags_index: StringIndexFilter
  tags: StringFilter
  and: [SignatureInfoFilter]
  or: [SignatureInfoFilter]
  nor: [SignatureInfoFilter]
  not: SignatureInfoFilter
}

input ContainerImageRepoTagElemMatchFilter {
  and: [ContainerImageRepoTagFilter]
  or: [ContainerImageRepoTagFilter]
  nor: [ContainerImageRepoTagFilter]
  not: ContainerImageRepoTagFilter
}

input ContainerImageRepoTagFilter {
  added_date: DateTimeFilter
  manifest_schema1_digest: StringFilter
  name: StringFilter
  removed_date: DateTimeFilter
  and: [ContainerImageRepoTagFilter]
  or: [ContainerImageRepoTagFilter]
  nor: [ContainerImageRepoTagFilter]
  not: ContainerImageRepoTagFilter
}

input SortBy {
  field: String
  order: SortDirectionEnum
}

enum SortDirectionEnum {
  ASC
  DESC
}

type ContainerImageRepoEdges {
  image_advisory: RedHatContainerAdvisoryResponse
  repository: ContainerRepositoryResponse
}

type RedHatContainerAdvisoryResponse {
  data: RedHatContainerAdvisory
  error: ResponseError
}

# Advisory associated with RH container image.
type RedHatContainerAdvisory {
  # MongoDB unique _id
  _id: String

  # The content type of advisory. i.e. for CONTAINER or RPM.
  content_type: String

  # Advisory description.
  description: String
  object_type: String

  # Severity of the advisory.
  severity: String

  # The date the image advisory shipped.
  ship_date: DateTime

  # The solution of the advisory.
  solution: String

  # Short summary of the advisory.
  synopsis: String

  # Topic of the advisory.
  topic: String

  # The type of advisory. i.e. for RHSA:2016-1001 the type is 'RHSA'.
  type: String

  # Array of CVEs fixed by this advisory.
  cves: [CVE]

  # Array of issues fixed by this advisory.
  issues: [Issue]

  # The date when the entry was created. Value is created automatically on creation.
  creation_date: DateTime

  # The date when the entry was last updated.
  last_update_date: DateTime
}

# CVE fixed by an advisory.
type CVE {
  # Unique identifier of the issue in the issue tracking system.
  id: String

  # Publicly accessible URL of the issue information.
  url: String
}

# Issue fixed by an advisory.
type Issue {
  # Unique identifier of the issue in the issue tracking system.
  id: String

  # Hostname of the issue tracking system used.
  issue_tracker: String

  # Publicly accessible URL of the issue information.
  url: String
}

type ContainerRepositoryResponse {
  data: ContainerRepository
  error: ResponseError
}

# Contains metadata associated with Red Hat and ISV repositories
type ContainerRepository {
  # The application categories (types).
  application_categories: [String]

  # Contains unique list of all container architectures for the given repository.
  architectures: [String]

  # Denote which tags to be used for auto-rebuilding processes.
  auto_rebuild_tags: [String]

  # Flag indicating whether the repository is still beta or not.
  beta: Boolean @deprecated(reason: "Deprecated")

  # What build categories does this fall into, such as standalone, s2i builder, etc.
  build_categories: [String]

  # Flag indicating whether the repository has opted-in to auto-release auto-built images.
  can_auto_release_cve_rebuild: Boolean
  cdn_base_url: String

  # To provide customers information which yum repos to enable ability to update the container content.
  content_sets: [String]
    @deprecated(reason: "Use containerImage.content_sets instead.")

  # Capture and provide an inventory of grades corresponding to the tags in the relevant contents stream.
  content_stream_grades: [RepositoryContentStreamGrades]

  # Capture and provide an inventory of tags corresponding to the content streams.
  content_stream_tags: [String]

  # Flag indicating whether the repository is still supported or not.
  deprecated: Boolean @deprecated(reason: "Deprecated")

  # Description of the repository.
  description: String
  display_data: RepositoryDisplayData

  # Links to marketing and doc collateral including categorization (solution brief, white paper, demo video, etc.) supposed to be displayed on the product page (NOT documentation tab on image overview tab).
  documentation_links: [RepositoryDocumentationLink]
  eol_date: DateTime

  # Date until the freshness grades for this repository are unknown.
  freshness_grades_unknown_until_date: DateTime

  # Defines  whether a repository contains multiple image streams.
  includes_multiple_content_streams: Boolean

  # Designates whether a repository is community-supported.
  is_community_supported: Boolean

  # ID of the project in for ISV repositories.
  isv_pid: String

  # Manually overriden label values.  When set, should be taken instead of label set on the image.
  label_override: RepositoryLabelOverride @deprecated(reason: "Deprecated")

  # Set of metrics about the repository.
  metrics: RepositoryMetrics

  # Namespace of the repository.
  namespace: String

  # Repository is intended for non-production use only.
  non_production_only: Boolean
  object_type: String

  # Indicates if images in this repository are allowed to run super-privileged.
  privileged_images_allowed: Boolean

  # ID of the project in PRM. Only for ISV repositories.
  prm_project_id: String

  # Reference to the product for this repository by id.
  product_id: String
    @deprecated(
      reason: "The product data has been replaced by product listings for ISVs. For RH products it will be replaced June 2020."
    )

  # List of unique identifiers for the product listings.
  product_listings: [String]

  # Map repositories to specific product versions.
  product_versions: [String]
  protected_for_pull: Boolean

  # Indicates whether the repository requires subscription or other access restrictions for search.
  protected_for_search: Boolean

  # Indicates that the repository does not have any images in it or has been deleted.
  published: Boolean

  # Hostname of the registry where the repository can be accessed.
  registry: String

  # Consumed by the Registry Proxy so that it can route users to the proper backend registry (e.g. Pulp or Quay).
  registry_target: String

  # The release categories of a repository.
  release_categories: [String]

  # Defines repository to point to in case this one is deprecated.
  replaced_by_repository_name: String

  # Combination of image repository and namespace.
  repository: String

  # Flag indicating whether (false) the repository is published on the legacy registry (registry.access.redhat.com), or (true) can only be published to registry.redhat.io.
  requires_terms: Boolean

  # Describes what the image can be run on.
  runs_on: RepositoryRunsOn

  # Flag indicating whether images associated with this repo are included in workflows where non-binary container images are published alongside their binary counterparts.
  source_container_image_enabled: Boolean

  # The support levels of a repository.
  support_levels: [String]

  # Flag indicating whether the repository is in tech preview or not.
  tech_preview: Boolean @deprecated(reason: "Deprecated")

  # Total size of all images in bytes.
  total_size_bytes: Int

  # Total size of all uncompressed images in bytes.
  total_uncompressed_size_bytes: Int

  # When populated this field will override the content on the 'get this image' tab in red hat container catalog.
  ui_get_this_image_override: String

  # Label of the vendor that owns this repository.
  vendor_label: String

  # MongoDB unique _id
  _id: String

  # The date when the entry was created. Value is created automatically on creation.
  creation_date: DateTime

  # The date when the entry was last updated.
  last_update_date: DateTime
  edges: ContainerRepositoryEdges
}

# Content stream grades associated with image tag.
type RepositoryContentStreamGrades {
  # Name of floating tag associated with the image.
  tag: String

  # Current image security grade.
  grade: String
}

# Display data for Catalog.
type RepositoryDisplayData {
  # The long description of the repository.
  long_description: String

  # Markdown for the long description of the repository.
  long_description_markdown: String

  # The repository name.
  name: String

  # The io_openshift_tags.
  openshift_tags: String

  # The short description of the repository.
  short_description: String
}

# Information about documentation link.
type RepositoryDocumentationLink {
  # Title of the link.
  title: String

  # The type of collateral.
  type: String

  # The URL for the documentation.
  url: String
}

# Label override data.
type RepositoryLabelOverride {
  # Override for 'description' label.
  description: String

  # Override for 'io.k8s.display-name' label.
  io_k8s_displayName: String

  # Override for 'io.openshift.tags' label.
  io_openshift_tags: String

  # Override for 'summary' label.
  summary: String
}

# Metrics information.
type RepositoryMetrics {
  # The date and time when these metrics were last updated for the repository.
  last_update_date: DateTime

  # The number of pulls in the last 30 days for the repository.
  pulls_in_last_30_days: Int
}

# Describes what the image can be run on.
type RepositoryRunsOn {
  # Can the image run on openshift_online.
  openshift_online: Boolean
}

type ContainerRepositoryEdges {
  certification_project(
    page_size: Int = 50
    page: Int = 0
    filter: CertificationProjectFilter
    sort_by: [SortBy]
  ): CertificationProjectPaginatedResponse
  images(
    page_size: Int = 50
    page: Int = 0
    filter: ContainerImageFilter
    sort_by: [SortBy]
  ): ContainerImagePaginatedResponse
  product_listings(
    page_size: Int = 50
    page: Int = 0
    filter: ProductListingFilter
    sort_by: [SortBy]
  ): ProductListingPaginatedResponse
  operator_bundles(
    page_size: Int = 50
    page: Int = 0
    filter: OperatorBundleFilter
    sort_by: [SortBy]
  ): OperatorBundlePaginatedResponse
  replaced_by_repository: ContainerRepositoryResponse
  vendor: ContainerVendorResponse
}

type CertificationProjectPaginatedResponse {
  data: [CertificationProject]
  error: ResponseError
  page: Int
  page_size: Int
  total: Int
}

# Certification project information.
type CertificationProject {
  # Certification Date.
  certification_date: DateTime

  # Certification Status.
  certification_status: String

  # Certification User.
  certification_user: Int

  # Contacts for certification project.
  contacts: [CertProjectContacts]
  container: CertProjectContainer

  # Configuration specific to Helm Chart projects.
  helm_chart: CertProjectHelmChart
  drupal: CertProjectDrupal @deprecated(reason: "Deprecated.")
  marketplace: CertProjectMarketplace

  # The owner provided name of the certification project.
  name: String

  # Operator Distribution.
  operator_distribution: String

  # Red Hat Org ID / account_id from Red Hat SSO. Also corresponds to company_org_id in Red Hat Connect.
  org_id: Int

  # Unique identifier for the product listing.
  product_listings: [String]

  # Status of the certification project.
  project_status: String

  # Who published the certification project.
  published_by: String
  redhat: CertProjectRedhat
  self_certification: CertProjectSelfCertification

  # Certification project type.
  type: String

  # MongoDB unique _id
  _id: String

  # The date when the entry was created. Value is created automatically on creation.
  creation_date: DateTime

  # The date when the entry was last updated.
  last_update_date: DateTime
  edges: CertificationProjectEdges
}

# Contact info.
type CertProjectContacts {
  email_address: Email
  type: String
}

scalar Email

# Container related information.
type CertProjectContainer {
  # The application categories (types).
  application_categories: [String]

  # Once a container is certified it is automatically published. Auto-publish must be enabled in order to set up automatic rebuilds. Auto-publish must always be enabled when auto-rebuilding is enabled.
  auto_publish: Boolean

  # Auto rebuild enabled.
  auto_rebuild: Boolean

  # Distribution approval obtained.
  distribution_approval: Boolean

  # Distribution method.
  distribution_method: String

  # ID of the project in for ISV repositories.
  isv_pid: String

  # Kubernetes objects for operator registry projects. Value has to be a valid YAML.
  kube_objects: OpenPGPEncrypted

  # Docker config for operator registry projects. Value has to be a valid JSON.
  docker_config_json: OpenPGPEncrypted

  # OS Content Type.
  os_content_type: String

  # Passed RH Cert.
  passed_rhcert: Boolean

  # A container needs to run in a privileged state.
  privileged: Boolean
  published: Boolean
    @deprecated(reason: "The field was replaced by certification_status.")

  # Hostname of the registry where the repository can be accessed.
  # Examples: registry.company.com assumes the default port, 443. registry.company.com:5000 repository path with optional port specified.
  # It is only applicable for projects with an 'external' distribution method.
  registry: String

  # Note: These instructions will be displayed in the Red Hat Container Catalog as is. Please modify the following template as it suits your needs.
  registry_override_instruct: String

  # Release category.
  release_category: String

  # Path to the container repository as found in the registry.
  #
  # Examples:
  # path/to/repository
  # repository
  #
  # This field can only be edited when there are no published containers.
  # It is only applicable for projects with an 'external' distribution method.
  repository: String

  # The repository description is displayed on the container
  # catalog repository overview page.
  repository_description: String

  # This should represent your product (or the component if your product consists of multiple containers)
  # and a major version. For example, you could use names like jboss-server7, or agent5.
  #
  # This value is only editable when there are no published containers in this project.
  # It is only applicable for projects that do not have the 'external' distribution method.
  repository_name: String

  # Service Account Secret.
  service_account_secret: String

  # Short description of the container.
  short_description: String

  # Supported Platforms.
  support_platforms: [String]

  # Container type.
  # Field is required, if project type is 'Container', and the field is immutable for Partners after creation.
  type: String

  # Filename other than the default Dockerfile or a path to a Dockerfile in a subdirectory.
  source_dockerfile: String

  # Force the build to ignore cached layers and rerun all steps of the Dockerfile.
  build_no_cache: Boolean

  # Override default location (root directory) for applications within a subdirectory.
  source_context_dir: String

  # Whether Red Hat will build your container.
  build_service: Boolean

  # The specific Git branch to checkout.
  source_ref: String

  # The URL to the source used for the build.
  # For example: 'https://github.com/openshift/ruby-hello-world
  source_uri: URI

  # Base64 encoded SSH private key in PEM format. Used to pull the source.
  source_ssh_private_key: Base64OpenPGPEncrypted

  # GitHub users authorized to submit a certification pull request.
  github_usernames: [String]
  edges: CertProjectContainerEdges
}

scalar OpenPGPEncrypted

scalar URI

scalar Base64OpenPGPEncrypted

type CertProjectContainerEdges {
  repository: ContainerRepositoryResponse
}

# Helm chart related information.
type CertProjectHelmChart {
  # How your Helm Chart is distributed.
  distribution_method: String

  # The Helm Chart name as it will appear in GitHub.
  chart_name: String

  # URL to the externally distributed Helm Chart repository. This is not used if the chart is distributed via Red Hat.
  repository: URI

  # Instructions for users to access an externally distributed Helm Chart.
  distribution_instructions: String

  # Base64 encoded PGP public key. Used to sign result submissions.
  public_pgp_key: String

  # URL to the user submitted github pull request for this project.
  github_pull_request: URI

  # Short description of the Helm Chart.
  short_description: String

  # Long description of the Helm Chart.
  long_description: String

  # The application categories (types).
  application_categories: [String]

  # GitHub users authorized to submit a certification pull request.
  github_usernames: [String]
}

# Drupal related information.
type CertProjectDrupal {
  # Company node ID from Red Hat Connect.
  company_id: Int

  # Relation ID for certification project.
  relation: Int

  # Zone for certification project.
  zone: String
}

# Marketplace related information.
type CertProjectMarketplace {
  enablement_status: String
  enablement_url: URI
  listing_url: URI
  published: Boolean
}

# Red Hat projects related information.
type CertProjectRedhat {
  # Red Hat Product ID.
  product_id: Int

  # Red Hat product name.
  product_name: String

  # Red Hat Product Version.
  product_version: String

  # Red Hat Product Version.
  product_version_id: Int
}

# Red Hat projects related information.
type CertProjectSelfCertification {
  # Application Profiler.
  app_profiler: Boolean

  # Application Runs on App Type.
  app_runs_on_app_type: Boolean

  # Whether the Self Certification Evidence URL requires a customer login.
  auth_login: Boolean

  # Self Certification Evidence URL.
  certification_url: URI

  # Can Commercially Support on App Type.
  comm_support_on_app_type: Boolean

  # Self Certification Requested.
  requested: Boolean

  # TsaNET Member.
  tsanet_member: Boolean
}

type CertificationProjectEdges {
  vendor: ContainerVendorResponse
  build_requests(
    page_size: Int = 50
    page: Int = 0
    filter: CertProjectBuildRequestFilter
    sort_by: [SortBy]
  ): CertProjectBuildRequestPaginatedResponse
  scan_requests(
    page_size: Int = 50
    page: Int = 0
    filter: CertProjectScanRequestFilter
    sort_by: [SortBy]
  ): CertProjectScanRequestPaginatedResponse
  tag_requests(
    page_size: Int = 50
    page: Int = 0
    filter: CertProjectTagRequestFilter
    sort_by: [SortBy]
  ): CertProjectTagRequestPaginatedResponse
}

type ContainerVendorResponse {
  data: ContainerVendor
  error: ResponseError
}

# Stores information about a Vendor
type ContainerVendor {
  # URL to the vendor's main website.
  company_url: URI

  # General contact information for the vendor, to be displayed on the vendor page on RHCC.
  contact: ContainerVendorContact
  description: String

  # Company node ID from Red Hat Connect.
  drupal_company_id: Int

  # The industry / vertical the vendor belongs to.
  industries: [String]
  label: String

  # A flag that determines if vendor label can be changed.
  label_locked: Boolean
  logo_url: URI
  name: String
  object_type: String

  # Indicate that the vendor has been published.
  published: Boolean
  registry_urls: [String]

  # RSS feed for vendor.
  rss_feed_url: URI

  # Token for outbound namespace for pulling published marketplace images.
  service_account_token: String
  social_media_links: [ContainerVendorSocialMediaLinks]

  # Red Hat Org ID / account_id from Red Hat SSO. Also corresponds to company_org_id in Red Hat Connect.
  org_id: Int

  # MongoDB unique _id
  _id: String

  # The date when the entry was created. Value is created automatically on creation.
  creation_date: DateTime

  # The date when the entry was last updated.
  last_update_date: DateTime
  edges: ContainerVendorEdges
}

# Contact information
type ContainerVendorContact {
  # General contact email address.
  email: String

  # General contact phone number.
  phone: String
}

# Social media links.
type ContainerVendorSocialMediaLinks {
  # The name of the social media provider.
  name: String

  # The URL to the social media site for the vendor.
  url: URI
}

type ContainerVendorEdges {
  product_listings_org_id(
    page_size: Int = 50
    page: Int = 0
    filter: ProductListingFilter
    sort_by: [SortBy]
  ): ProductListingPaginatedResponse
  product_listings_label(
    page_size: Int = 50
    page: Int = 0
    filter: ProductListingFilter
    sort_by: [SortBy]
  ): ProductListingPaginatedResponse
  repositories(
    page_size: Int = 50
    page: Int = 0
    filter: ContainerRepositoryFilter
    sort_by: [SortBy]
  ): ContainerRepositoryPaginatedResponse
}

type ProductListingPaginatedResponse {
  data: [ProductListing]
  error: ResponseError
  page: Int
  page_size: Int
  total: Int
}

# Product listings define a marketing page in the Ecosystem Catalog. It allows you to group repos and showcase what they accomplish together as an application. In the case of operators, your CSV file populates OperatorHub, which can only be viewed in cluster through OpenShift. Your product listing is publicly visible in the Ecosystem Catalog so anyone can know that it is offered.
type ProductListing {
  badges: [BadgesItems] @deprecated(reason: "Deprecated")
  category: String

  # List of unique identifiers for the certification project.
  cert_projects: [String]

  # This field is required when the product listing is published.
  contacts: [ContactsItems]

  # This field is required when the product listing is published.
  descriptions: Descriptions

  # Company node ID from Red Hat Connect. Read only.
  drupal_company_id: Int

  # This field is required when the product listing is published.
  faqs: [FAQSItems]

  # This field is required when the product listing is published.
  features: [FeaturesItems]

  # This field is required when the product listing is published.
  functional_categories: [String]
  legal: Legal

  # This field is required when the product listing is published.
  linked_resources: [LinkedResourcesItems]
  logo: Logo
  marketplace: Marketplace
    @deprecated(reason: "This field has been moved to certProject.")
  name: String
  published: Boolean

  # Flag determining if product listing is considered to be deleted. Product listing can be deleted only if it is not published. Value is set to False by default.
  deleted: Boolean
  quick_start_configuration: QuickStartConfiguration

  # List of unique identifiers for the repository.
  repositories: [String]

  # This field is required when the product listing is published.
  search_aliases: [SearchAliasesItems]
  support: Support
  type: String
  vendor_label: String
  operator_bundles: [OperatorBundlesItems]

  # Red Hat Org ID / account_id from Red Hat SSO. Also corresponds to company_org_id in Red Hat Connect.
  org_id: Int

  # MongoDB unique _id
  _id: String

  # The date when the entry was created. Value is created automatically on creation.
  creation_date: DateTime

  # The date when the entry was last updated.
  last_update_date: DateTime
  edges: ProductListingEdges
}

#
type BadgesItems {
  badge: String
  project_id: ObjectID
}

scalar ObjectID

#
type ContactsItems {
  email_address: Email
  type: String
}

# This field is required when the product listing is published.
type Descriptions {
  long: String
  short: String
}

# This field is required when the product listing is published.
type FAQSItems {
  answer: String
  question: String
}

# This field is required when the product listing is published.
type FeaturesItems {
  description: String
  title: String
}

#
type Legal {
  description: String
  license_agreement_url: URI
  privacy_policy_url: URI
}

#
type LinkedResourcesItems {
  category: String
  description: String
  thumbnail_url: URI
  title: String
  type: String
  url: URI
}

#
type Logo {
  url: URI
}

#
type Marketplace {
  enablement_status: String
  enablement_url: URI
  listing_url: URI
  published: Boolean
}

#
type QuickStartConfiguration {
  instructions: String
}

#
type SearchAliasesItems {
  key: String
  value: String
}

# This field is required when the product listing is published.
type Support {
  description: String
  email_address: Email
  phone_number: String
  url: URI
}

#
type OperatorBundlesItems {
  # Bundle unique identifier
  _id: ObjectID

  # Bundle package name
  package: String
  capabilities: [String]
}

type ProductListingEdges {
  vendor: ContainerVendorResponse
}

input ProductListingFilter {
  badges_size: ListSizeFilter
  badges_index: StringIndexFilter
  badges_elemMatch: BadgesItemsElemMatchFilter
  badges: BadgesItemsFilter
  category: StringFilter
  cert_projects_size: ListSizeFilter
  cert_projects_index: StringIndexFilter
  cert_projects: StringFilter
  contacts_size: ListSizeFilter
  contacts_index: StringIndexFilter
  contacts_elemMatch: ContactsItemsElemMatchFilter
  contacts: ContactsItemsFilter
  descriptions: DescriptionsFilter
  drupal_company_id: IntFilter
  faqs_size: ListSizeFilter
  faqs_index: StringIndexFilter
  faqs_elemMatch: FAQSItemsElemMatchFilter
  faqs: FAQSItemsFilter
  features_size: ListSizeFilter
  features_index: StringIndexFilter
  features_elemMatch: FeaturesItemsElemMatchFilter
  features: FeaturesItemsFilter
  functional_categories_size: ListSizeFilter
  functional_categories_index: StringIndexFilter
  functional_categories: StringFilter
  legal: LegalFilter
  linked_resources_size: ListSizeFilter
  linked_resources_index: StringIndexFilter
  linked_resources_elemMatch: LinkedResourcesItemsElemMatchFilter
  linked_resources: LinkedResourcesItemsFilter
  logo: LogoFilter
  marketplace: MarketplaceFilter
  name: StringFilter
  published: BooleanFilter
  deleted: BooleanFilter
  quick_start_configuration: QuickStartConfigurationFilter
  repositories_size: ListSizeFilter
  repositories_index: StringIndexFilter
  repositories: StringFilter
  search_aliases_size: ListSizeFilter
  search_aliases_index: StringIndexFilter
  search_aliases_elemMatch: SearchAliasesItemsElemMatchFilter
  search_aliases: SearchAliasesItemsFilter
  support: SupportFilter
  type: StringFilter
  vendor_label: StringFilter
  operator_bundles_size: ListSizeFilter
  operator_bundles_index: StringIndexFilter
  operator_bundles_elemMatch: OperatorBundlesItemsElemMatchFilter
  operator_bundles: OperatorBundlesItemsFilter
  org_id: IntFilter
  _id: StringFilter
  creation_date: DateTimeFilter
  last_update_date: DateTimeFilter
  and: [ProductListingFilter]
  or: [ProductListingFilter]
  nor: [ProductListingFilter]
  not: ProductListingFilter
}

input BadgesItemsElemMatchFilter {
  and: [BadgesItemsFilter]
  or: [BadgesItemsFilter]
  nor: [BadgesItemsFilter]
  not: BadgesItemsFilter
}

input BadgesItemsFilter {
  badge: StringFilter
  project_id: StringFilter
  and: [BadgesItemsFilter]
  or: [BadgesItemsFilter]
  nor: [BadgesItemsFilter]
  not: BadgesItemsFilter
}

input ContactsItemsElemMatchFilter {
  and: [ContactsItemsFilter]
  or: [ContactsItemsFilter]
  nor: [ContactsItemsFilter]
  not: ContactsItemsFilter
}

input ContactsItemsFilter {
  email_address: EmailFilter
  type: StringFilter
  and: [ContactsItemsFilter]
  or: [ContactsItemsFilter]
  nor: [ContactsItemsFilter]
  not: ContactsItemsFilter
}

input EmailFilter {
  eq: Email
  ne: Email
  le: Email
  lt: Email
  ge: Email
  gt: Email
  in: [Email]
}

input DescriptionsFilter {
  long: StringFilter
  short: StringFilter
  and: [DescriptionsFilter]
  or: [DescriptionsFilter]
  nor: [DescriptionsFilter]
  not: DescriptionsFilter
}

input FAQSItemsElemMatchFilter {
  and: [FAQSItemsFilter]
  or: [FAQSItemsFilter]
  nor: [FAQSItemsFilter]
  not: FAQSItemsFilter
}

input FAQSItemsFilter {
  answer: StringFilter
  question: StringFilter
  and: [FAQSItemsFilter]
  or: [FAQSItemsFilter]
  nor: [FAQSItemsFilter]
  not: FAQSItemsFilter
}

input FeaturesItemsElemMatchFilter {
  and: [FeaturesItemsFilter]
  or: [FeaturesItemsFilter]
  nor: [FeaturesItemsFilter]
  not: FeaturesItemsFilter
}

input FeaturesItemsFilter {
  description: StringFilter
  title: StringFilter
  and: [FeaturesItemsFilter]
  or: [FeaturesItemsFilter]
  nor: [FeaturesItemsFilter]
  not: FeaturesItemsFilter
}

input LegalFilter {
  description: StringFilter
  license_agreement_url: URIFilter
  privacy_policy_url: URIFilter
  and: [LegalFilter]
  or: [LegalFilter]
  nor: [LegalFilter]
  not: LegalFilter
}

input URIFilter {
  eq: URI
  ne: URI
  le: URI
  lt: URI
  ge: URI
  gt: URI
  in: [URI]
}

input LinkedResourcesItemsElemMatchFilter {
  and: [LinkedResourcesItemsFilter]
  or: [LinkedResourcesItemsFilter]
  nor: [LinkedResourcesItemsFilter]
  not: LinkedResourcesItemsFilter
}

input LinkedResourcesItemsFilter {
  category: StringFilter
  description: StringFilter
  thumbnail_url: URIFilter
  title: StringFilter
  type: StringFilter
  url: URIFilter
  and: [LinkedResourcesItemsFilter]
  or: [LinkedResourcesItemsFilter]
  nor: [LinkedResourcesItemsFilter]
  not: LinkedResourcesItemsFilter
}

input LogoFilter {
  url: URIFilter
  and: [LogoFilter]
  or: [LogoFilter]
  nor: [LogoFilter]
  not: LogoFilter
}

input MarketplaceFilter {
  enablement_status: StringFilter
  enablement_url: URIFilter
  listing_url: URIFilter
  published: BooleanFilter
  and: [MarketplaceFilter]
  or: [MarketplaceFilter]
  nor: [MarketplaceFilter]
  not: MarketplaceFilter
}

input QuickStartConfigurationFilter {
  instructions: StringFilter
  and: [QuickStartConfigurationFilter]
  or: [QuickStartConfigurationFilter]
  nor: [QuickStartConfigurationFilter]
  not: QuickStartConfigurationFilter
}

input SearchAliasesItemsElemMatchFilter {
  and: [SearchAliasesItemsFilter]
  or: [SearchAliasesItemsFilter]
  nor: [SearchAliasesItemsFilter]
  not: SearchAliasesItemsFilter
}

input SearchAliasesItemsFilter {
  key: StringFilter
  value: StringFilter
  and: [SearchAliasesItemsFilter]
  or: [SearchAliasesItemsFilter]
  nor: [SearchAliasesItemsFilter]
  not: SearchAliasesItemsFilter
}

input SupportFilter {
  description: StringFilter
  email_address: EmailFilter
  phone_number: StringFilter
  url: URIFilter
  and: [SupportFilter]
  or: [SupportFilter]
  nor: [SupportFilter]
  not: SupportFilter
}

input OperatorBundlesItemsElemMatchFilter {
  and: [OperatorBundlesItemsFilter]
  or: [OperatorBundlesItemsFilter]
  nor: [OperatorBundlesItemsFilter]
  not: OperatorBundlesItemsFilter
}

input OperatorBundlesItemsFilter {
  _id: StringFilter
  package: StringFilter
  capabilities_size: ListSizeFilter
  capabilities_index: StringIndexFilter
  capabilities: StringFilter
  and: [OperatorBundlesItemsFilter]
  or: [OperatorBundlesItemsFilter]
  nor: [OperatorBundlesItemsFilter]
  not: OperatorBundlesItemsFilter
}

type ContainerRepositoryPaginatedResponse {
  data: [ContainerRepository]
  error: ResponseError
  page: Int
  page_size: Int
  total: Int
}

input ContainerRepositoryFilter {
  application_categories_size: ListSizeFilter
  application_categories_index: StringIndexFilter
  application_categories: StringFilter
  architectures_size: ListSizeFilter
  architectures_index: StringIndexFilter
  architectures: StringFilter
  auto_rebuild_tags_size: ListSizeFilter
  auto_rebuild_tags_index: StringIndexFilter
  auto_rebuild_tags: StringFilter
  beta: BooleanFilter
  build_categories_size: ListSizeFilter
  build_categories_index: StringIndexFilter
  build_categories: StringFilter
  can_auto_release_cve_rebuild: BooleanFilter
  cdn_base_url: StringFilter
  content_sets_size: ListSizeFilter
  content_sets_index: StringIndexFilter
  content_sets: StringFilter
  content_stream_grades_size: ListSizeFilter
  content_stream_grades_index: StringIndexFilter
  content_stream_grades_elemMatch: RepositoryContentStreamGradesElemMatchFilter
  content_stream_grades: RepositoryContentStreamGradesFilter
  content_stream_tags_size: ListSizeFilter
  content_stream_tags_index: StringIndexFilter
  content_stream_tags: StringFilter
  deprecated: BooleanFilter
  description: StringFilter
  display_data: RepositoryDisplayDataFilter
  documentation_links_size: ListSizeFilter
  documentation_links_index: StringIndexFilter
  documentation_links_elemMatch: RepositoryDocumentationLinkElemMatchFilter
  documentation_links: RepositoryDocumentationLinkFilter
  eol_date: DateTimeFilter
  freshness_grades_unknown_until_date: DateTimeFilter
  includes_multiple_content_streams: BooleanFilter
  is_community_supported: BooleanFilter
  isv_pid: StringFilter
  label_override: RepositoryLabelOverrideFilter
  metrics: RepositoryMetricsFilter
  namespace: StringFilter
  non_production_only: BooleanFilter
  object_type: StringFilter
  privileged_images_allowed: BooleanFilter
  prm_project_id: StringFilter
  product_id: StringFilter
  product_listings_size: ListSizeFilter
  product_listings_index: StringIndexFilter
  product_listings: StringFilter
  product_versions_size: ListSizeFilter
  product_versions_index: StringIndexFilter
  product_versions: StringFilter
  protected_for_pull: BooleanFilter
  protected_for_search: BooleanFilter
  published: BooleanFilter
  registry: StringFilter
  registry_target: StringFilter
  release_categories_size: ListSizeFilter
  release_categories_index: StringIndexFilter
  release_categories: StringFilter
  replaced_by_repository_name: StringFilter
  repository: StringFilter
  requires_terms: BooleanFilter
  runs_on: RepositoryRunsOnFilter
  source_container_image_enabled: BooleanFilter
  support_levels_size: ListSizeFilter
  support_levels_index: StringIndexFilter
  support_levels: StringFilter
  tech_preview: BooleanFilter
  total_size_bytes: IntFilter
  total_uncompressed_size_bytes: IntFilter
  ui_get_this_image_override: StringFilter
  vendor_label: StringFilter
  _id: StringFilter
  creation_date: DateTimeFilter
  last_update_date: DateTimeFilter
  and: [ContainerRepositoryFilter]
  or: [ContainerRepositoryFilter]
  nor: [ContainerRepositoryFilter]
  not: ContainerRepositoryFilter
}

input RepositoryContentStreamGradesElemMatchFilter {
  and: [RepositoryContentStreamGradesFilter]
  or: [RepositoryContentStreamGradesFilter]
  nor: [RepositoryContentStreamGradesFilter]
  not: RepositoryContentStreamGradesFilter
}

input RepositoryContentStreamGradesFilter {
  tag: StringFilter
  grade: StringFilter
  and: [RepositoryContentStreamGradesFilter]
  or: [RepositoryContentStreamGradesFilter]
  nor: [RepositoryContentStreamGradesFilter]
  not: RepositoryContentStreamGradesFilter
}

input RepositoryDisplayDataFilter {
  long_description: StringFilter
  long_description_markdown: StringFilter
  name: StringFilter
  openshift_tags: StringFilter
  short_description: StringFilter
  and: [RepositoryDisplayDataFilter]
  or: [RepositoryDisplayDataFilter]
  nor: [RepositoryDisplayDataFilter]
  not: RepositoryDisplayDataFilter
}

input RepositoryDocumentationLinkElemMatchFilter {
  and: [RepositoryDocumentationLinkFilter]
  or: [RepositoryDocumentationLinkFilter]
  nor: [RepositoryDocumentationLinkFilter]
  not: RepositoryDocumentationLinkFilter
}

input RepositoryDocumentationLinkFilter {
  title: StringFilter
  type: StringFilter
  url: StringFilter
  and: [RepositoryDocumentationLinkFilter]
  or: [RepositoryDocumentationLinkFilter]
  nor: [RepositoryDocumentationLinkFilter]
  not: RepositoryDocumentationLinkFilter
}

input RepositoryLabelOverrideFilter {
  description: StringFilter
  io_k8s_displayName: StringFilter
  io_openshift_tags: StringFilter
  summary: StringFilter
  and: [RepositoryLabelOverrideFilter]
  or: [RepositoryLabelOverrideFilter]
  nor: [RepositoryLabelOverrideFilter]
  not: RepositoryLabelOverrideFilter
}

input RepositoryMetricsFilter {
  last_update_date: DateTimeFilter
  pulls_in_last_30_days: IntFilter
  and: [RepositoryMetricsFilter]
  or: [RepositoryMetricsFilter]
  nor: [RepositoryMetricsFilter]
  not: RepositoryMetricsFilter
}

input RepositoryRunsOnFilter {
  openshift_online: BooleanFilter
  and: [RepositoryRunsOnFilter]
  or: [RepositoryRunsOnFilter]
  nor: [RepositoryRunsOnFilter]
  not: RepositoryRunsOnFilter
}

type CertProjectBuildRequestPaginatedResponse {
  data: [CertProjectBuildRequest]
  error: ResponseError
  page: Int
  page_size: Int
  total: Int
}

# Contain status and related metadata of a certProject build request.
type CertProjectBuildRequest {
  # Unique identifier for the certification project.
  cert_project: ObjectID

  # The request status
  status: String

  # The tag that the container image gets when build is done.
  tag: String

  # MongoDB unique _id
  _id: String

  # The date when the entry was created. Value is created automatically on creation.
  creation_date: DateTime

  # The date when the entry was last updated.
  last_update_date: DateTime

  # Red Hat Org ID / account_id from Red Hat SSO. Also corresponds to company_org_id in Red Hat Connect.
  org_id: Int

  # An explanatory message to a request status.
  status_message: String
  edges: CertProjectBuildRequestEdges
}

type CertProjectBuildRequestEdges {
  cert_project: CertificationProjectResponse
  logs: CertProjectBuildLogResponse
}

type CertificationProjectResponse {
  data: CertificationProject
  error: ResponseError
}

type CertProjectBuildLogResponse {
  data: CertProjectBuildLog
  error: ResponseError
}

# Contain a certification project build's logs.
type CertProjectBuildLog {
  # Retrieved log for a certification project build.
  log: String
}

input CertProjectBuildRequestFilter {
  cert_project: StringFilter
  status: StringFilter
  tag: StringFilter
  _id: StringFilter
  creation_date: DateTimeFilter
  last_update_date: DateTimeFilter
  org_id: IntFilter
  status_message: StringFilter
  and: [CertProjectBuildRequestFilter]
  or: [CertProjectBuildRequestFilter]
  nor: [CertProjectBuildRequestFilter]
  not: CertProjectBuildRequestFilter
}

type CertProjectScanRequestPaginatedResponse {
  data: [CertProjectScanRequest]
  error: ResponseError
  page: Int
  page_size: Int
  total: Int
}

# Contain status and related metadata of a certProject scan request.
type CertProjectScanRequest {
  # Unique identifier for the certification project.
  cert_project: ObjectID

  # URL pointing to the location of DCI logs.
  external_tests_link: URI

  # Image pull specification in repo@sha256:digest format.
  pull_spec: String

  # Unique identifier of an ISV certification scan
  scan_uuid: String

  # Container image tag associated with the scan request.
  tag: String

  # MongoDB unique _id
  _id: String

  # The date when the entry was created. Value is created automatically on creation.
  creation_date: DateTime

  # The date when the entry was last updated.
  last_update_date: DateTime

  # Red Hat Org ID / account_id from Red Hat SSO. Also corresponds to company_org_id in Red Hat Connect.
  org_id: Int

  # The request status
  status: String

  # An explanatory message to a request status.
  status_message: String
  edges: CertProjectScanRequestEdges
}

type CertProjectScanRequestEdges {
  cert_project: CertificationProjectResponse
}

input CertProjectScanRequestFilter {
  cert_project: StringFilter
  external_tests_link: URIFilter
  pull_spec: StringFilter
  scan_uuid: StringFilter
  tag: StringFilter
  _id: StringFilter
  creation_date: DateTimeFilter
  last_update_date: DateTimeFilter
  org_id: IntFilter
  status: StringFilter
  status_message: StringFilter
  and: [CertProjectScanRequestFilter]
  or: [CertProjectScanRequestFilter]
  nor: [CertProjectScanRequestFilter]
  not: CertProjectScanRequestFilter
}

type CertProjectTagRequestPaginatedResponse {
  data: [CertProjectTagRequest]
  error: ResponseError
  page: Int
  page_size: Int
  total: Int
}

# Contain status and related metadata of a certProject tag request.
type CertProjectTagRequest {
  # Unique identifier for the certification project.
  cert_project: ObjectID

  # Container image id associated with the tag request.
  image_id: ObjectID

  # Operation performed during the tag request, e.g. publish
  operation: String

  # Container image tag associated with the tag request.
  tag: String

  # MongoDB unique _id
  _id: String

  # The date when the entry was created. Value is created automatically on creation.
  creation_date: DateTime

  # The date when the entry was last updated.
  last_update_date: DateTime

  # Red Hat Org ID / account_id from Red Hat SSO. Also corresponds to company_org_id in Red Hat Connect.
  org_id: Int

  # The request status
  status: String

  # An explanatory message to a request status.
  status_message: String
  edges: CertProjectTagRequestEdges
}

type CertProjectTagRequestEdges {
  cert_project: CertificationProjectResponse
  image: ContainerImageResponse
}

input CertProjectTagRequestFilter {
  cert_project: StringFilter
  image_id: StringFilter
  operation: StringFilter
  tag: StringFilter
  _id: StringFilter
  creation_date: DateTimeFilter
  last_update_date: DateTimeFilter
  org_id: IntFilter
  status: StringFilter
  status_message: StringFilter
  and: [CertProjectTagRequestFilter]
  or: [CertProjectTagRequestFilter]
  nor: [CertProjectTagRequestFilter]
  not: CertProjectTagRequestFilter
}

input CertificationProjectFilter {
  certification_date: DateTimeFilter
  certification_status: StringFilter
  certification_user: IntFilter
  contacts_size: ListSizeFilter
  contacts_index: StringIndexFilter
  contacts_elemMatch: CertProjectContactsElemMatchFilter
  contacts: CertProjectContactsFilter
  container: CertProjectContainerFilter
  helm_chart: CertProjectHelmChartFilter
  drupal: CertProjectDrupalFilter
  marketplace: CertProjectMarketplaceFilter
  name: StringFilter
  operator_distribution: StringFilter
  org_id: IntFilter
  product_listings_size: ListSizeFilter
  product_listings_index: StringIndexFilter
  product_listings: StringFilter
  project_status: StringFilter
  published_by: StringFilter
  redhat: CertProjectRedhatFilter
  self_certification: CertProjectSelfCertificationFilter
  type: StringFilter
  _id: StringFilter
  creation_date: DateTimeFilter
  last_update_date: DateTimeFilter
  and: [CertificationProjectFilter]
  or: [CertificationProjectFilter]
  nor: [CertificationProjectFilter]
  not: CertificationProjectFilter
}

input CertProjectContactsElemMatchFilter {
  and: [CertProjectContactsFilter]
  or: [CertProjectContactsFilter]
  nor: [CertProjectContactsFilter]
  not: CertProjectContactsFilter
}

input CertProjectContactsFilter {
  email_address: EmailFilter
  type: StringFilter
  and: [CertProjectContactsFilter]
  or: [CertProjectContactsFilter]
  nor: [CertProjectContactsFilter]
  not: CertProjectContactsFilter
}

input CertProjectContainerFilter {
  application_categories_size: ListSizeFilter
  application_categories_index: StringIndexFilter
  application_categories: StringFilter
  auto_publish: BooleanFilter
  auto_rebuild: BooleanFilter
  distribution_approval: BooleanFilter
  distribution_method: StringFilter
  isv_pid: StringFilter
  kube_objects: StringFilter
  docker_config_json: StringFilter
  os_content_type: StringFilter
  passed_rhcert: BooleanFilter
  privileged: BooleanFilter
  published: BooleanFilter
  registry: StringFilter
  registry_override_instruct: StringFilter
  release_category: StringFilter
  repository: StringFilter
  repository_description: StringFilter
  repository_name: StringFilter
  service_account_secret: StringFilter
  short_description: StringFilter
  support_platforms_size: ListSizeFilter
  support_platforms_index: StringIndexFilter
  support_platforms: StringFilter
  type: StringFilter
  source_dockerfile: StringFilter
  build_no_cache: BooleanFilter
  source_context_dir: StringFilter
  build_service: BooleanFilter
  source_ref: StringFilter
  source_uri: URIFilter
  source_ssh_private_key: StringFilter
  github_usernames_size: ListSizeFilter
  github_usernames_index: StringIndexFilter
  github_usernames: StringFilter
  and: [CertProjectContainerFilter]
  or: [CertProjectContainerFilter]
  nor: [CertProjectContainerFilter]
  not: CertProjectContainerFilter
}

input CertProjectHelmChartFilter {
  distribution_method: StringFilter
  chart_name: StringFilter
  repository: URIFilter
  distribution_instructions: StringFilter
  public_pgp_key: StringFilter
  github_pull_request: URIFilter
  short_description: StringFilter
  long_description: StringFilter
  application_categories_size: ListSizeFilter
  application_categories_index: StringIndexFilter
  application_categories: StringFilter
  github_usernames_size: ListSizeFilter
  github_usernames_index: StringIndexFilter
  github_usernames: StringFilter
  and: [CertProjectHelmChartFilter]
  or: [CertProjectHelmChartFilter]
  nor: [CertProjectHelmChartFilter]
  not: CertProjectHelmChartFilter
}

input CertProjectDrupalFilter {
  company_id: IntFilter
  relation: IntFilter
  zone: StringFilter
  and: [CertProjectDrupalFilter]
  or: [CertProjectDrupalFilter]
  nor: [CertProjectDrupalFilter]
  not: CertProjectDrupalFilter
}

input CertProjectMarketplaceFilter {
  enablement_status: StringFilter
  enablement_url: URIFilter
  listing_url: URIFilter
  published: BooleanFilter
  and: [CertProjectMarketplaceFilter]
  or: [CertProjectMarketplaceFilter]
  nor: [CertProjectMarketplaceFilter]
  not: CertProjectMarketplaceFilter
}

input CertProjectRedhatFilter {
  product_id: IntFilter
  product_name: StringFilter
  product_version: StringFilter
  product_version_id: IntFilter
  and: [CertProjectRedhatFilter]
  or: [CertProjectRedhatFilter]
  nor: [CertProjectRedhatFilter]
  not: CertProjectRedhatFilter
}

input CertProjectSelfCertificationFilter {
  app_profiler: BooleanFilter
  app_runs_on_app_type: BooleanFilter
  auth_login: BooleanFilter
  certification_url: URIFilter
  comm_support_on_app_type: BooleanFilter
  requested: BooleanFilter
  tsanet_member: BooleanFilter
  and: [CertProjectSelfCertificationFilter]
  or: [CertProjectSelfCertificationFilter]
  nor: [CertProjectSelfCertificationFilter]
  not: CertProjectSelfCertificationFilter
}

type OperatorBundlePaginatedResponse {
  data: [OperatorBundle]
  error: ResponseError
  page: Int
  page_size: Int
  total: Int
}

# An Operator Bundle is a container image that stores the Kubernetes manifests and metadata associated with an operator. A bundle is meant to represent a specific version of an operator.
type OperatorBundle {
  # Specific information from the CSV requested by customer.
  alm_examples: [ALMExample]

  # A subset of the \"metadata.annotations\" object from the CSV. Any annotations that are in the \"operators.openshift.io\" namespace that are strings of JSON will be expanded if it is valid JSON. Namespaces are not preserved due to their usage of special characters and all dashes are converted to underscores for consistency with other fields.
  annotations: OperatorBundleAnnotation

  # List of all supported architectures. An empty list means that architectures are unknown as CSV doesn't include architecture labels.
  architectures: [String]

  # The 'bundle' is the operator representation of a version of operator metadata. There is one bundle per operator/package channel for a particular OpenShift version.
  bundle: String

  # Pullspec of the operator bundle e.g. quay.io/foo/bar@sha256:digest.
  bundle_path: String

  # Digest from the bundle_path.
  bundle_path_digest: String

  # Specific information from the CSV requested by customer.
  capabilities: [String]

  # The channel for which this bundle is being released, e.g. \"amq-streams-1.5.x.
  channel_name: String

  # Full name of the package. Usually consists of package_name.version.
  csv_name: String

  # A public name to identify the Operator.
  csv_display_name: String

  # A thorough description of the Operator’s functionality in form of a markdown blob.
  csv_description: String

  # Short description of the CRD and operator functionality.
  csv_metadata_description: String

  # Indicate if the bundle is in an index image.
  in_index_img: Boolean

  # Specific information from the CSV requested by customer.
  install_modes: [OperatorBundleInstallMode]

  # If true then the channel is the default for this package, false otherwise.
  is_default_channel: Boolean

  # Indicate that the bundle is the latest version of a package in a channel for its associated OCP version (index image).
  latest_in_channel: Boolean

  # Specific OCP version for this bundle, e.g. \"4.5\".
  ocp_version: SemVer

  # Organization as understood by iib, e.g. \"redhat-marketplace\".
  organization: String

  # The name of the operator, e.g. \"amq-streams\".
  package: String

  # Specific information from the CSV requested by customer. Should correspond with values from alm_examples.
  provided_apis: [ProvidedAPIsItems]

  # Specific information from the CSV requested by customer.
  related_images: [RelatedImagesItems]

  # Where this bundle was collected from, e.g. \"quay.io/foo/bar:v4.5\".
  source_index_container_path: String

  # The operator version for this bundle
  version: SemVer

  # Original version of the bundle, used to recognize semver validity.
  version_original: String

  # Name of operator which the bundle replaces.
  replaces: String

  # List of skipped updates. See OLM upgrades documentation for more details.
  skips: [String]

  # String describing skipped versions.
  skip_range: String

  # MongoDB unique _id
  _id: String

  # The date when the entry was created. Value is created automatically on creation.
  creation_date: DateTime

  # The date when the entry was last updated.
  last_update_date: DateTime
}

# Specific information from the CSV requested by customer.
type ALMExample {
  #
  api_version: String

  #
  kind: String
}

# A subset of the \"metadata.annotations\" object from the CSV. Any annotations that are in the \"operators.openshift.io\" namespace that are strings of JSON will be expanded if it is valid JSON. Namespaces are not preserved due to their usage of special characters and all dashes are converted to underscores for consistency with other fields.
type OperatorBundleAnnotation {
  # The deserialized value of operators.openshift.io/infrastructure-features. This defaults to an empty array.
  infrastructure_features: [String]

  # The deserialized value of operators.openshift.io/valid-subscription. This defaults to an empty array.
  valid_subscription: [String]
}

# Specific information from the CSV requested by customer.
type OperatorBundleInstallMode {
  supported: Boolean
  type: String
}

scalar SemVer

# Specific information from the CSV requested by customer.Should correspond with values from alm_examples.
type ProvidedAPIsItems {
  group: String
  kind: String
  plural: String
  version: String
}

# Specific information from the CSV requested by customer.
type RelatedImagesItems {
  digest: String
  image: String
  name: String
}

input OperatorBundleFilter {
  alm_examples_size: ListSizeFilter
  alm_examples_index: StringIndexFilter
  alm_examples_elemMatch: ALMExampleElemMatchFilter
  alm_examples: ALMExampleFilter
  annotations: OperatorBundleAnnotationFilter
  architectures_size: ListSizeFilter
  architectures_index: StringIndexFilter
  architectures: StringFilter
  bundle: StringFilter
  bundle_path: StringFilter
  bundle_path_digest: StringFilter
  capabilities_size: ListSizeFilter
  capabilities_index: StringIndexFilter
  capabilities: StringFilter
  channel_name: StringFilter
  csv_name: StringFilter
  csv_display_name: StringFilter
  csv_description: StringFilter
  csv_metadata_description: StringFilter
  in_index_img: BooleanFilter
  install_modes_size: ListSizeFilter
  install_modes_index: StringIndexFilter
  install_modes_elemMatch: OperatorBundleInstallModeElemMatchFilter
  install_modes: OperatorBundleInstallModeFilter
  is_default_channel: BooleanFilter
  latest_in_channel: BooleanFilter
  ocp_version: StringFilter
  organization: StringFilter
  package: StringFilter
  provided_apis_size: ListSizeFilter
  provided_apis_index: StringIndexFilter
  provided_apis_elemMatch: ProvidedAPIsItemsElemMatchFilter
  provided_apis: ProvidedAPIsItemsFilter
  related_images_size: ListSizeFilter
  related_images_index: StringIndexFilter
  related_images_elemMatch: RelatedImagesItemsElemMatchFilter
  related_images: RelatedImagesItemsFilter
  source_index_container_path: StringFilter
  version: StringFilter
  version_original: StringFilter
  replaces: StringFilter
  skips_size: ListSizeFilter
  skips_index: StringIndexFilter
  skips: StringFilter
  skip_range: StringFilter
  _id: StringFilter
  creation_date: DateTimeFilter
  last_update_date: DateTimeFilter
  and: [OperatorBundleFilter]
  or: [OperatorBundleFilter]
  nor: [OperatorBundleFilter]
  not: OperatorBundleFilter
}

input ALMExampleElemMatchFilter {
  and: [ALMExampleFilter]
  or: [ALMExampleFilter]
  nor: [ALMExampleFilter]
  not: ALMExampleFilter
}

input ALMExampleFilter {
  api_version: StringFilter
  kind: StringFilter
  metadata: StringFilter
  spec: StringFilter
  and: [ALMExampleFilter]
  or: [ALMExampleFilter]
  nor: [ALMExampleFilter]
  not: ALMExampleFilter
}

input OperatorBundleAnnotationFilter {
  infrastructure_features_size: ListSizeFilter
  infrastructure_features_index: StringIndexFilter
  infrastructure_features: StringFilter
  valid_subscription_size: ListSizeFilter
  valid_subscription_index: StringIndexFilter
  valid_subscription: StringFilter
  and: [OperatorBundleAnnotationFilter]
  or: [OperatorBundleAnnotationFilter]
  nor: [OperatorBundleAnnotationFilter]
  not: OperatorBundleAnnotationFilter
}

input OperatorBundleInstallModeElemMatchFilter {
  and: [OperatorBundleInstallModeFilter]
  or: [OperatorBundleInstallModeFilter]
  nor: [OperatorBundleInstallModeFilter]
  not: OperatorBundleInstallModeFilter
}

input OperatorBundleInstallModeFilter {
  supported: BooleanFilter
  type: StringFilter
  and: [OperatorBundleInstallModeFilter]
  or: [OperatorBundleInstallModeFilter]
  nor: [OperatorBundleInstallModeFilter]
  not: OperatorBundleInstallModeFilter
}

input ProvidedAPIsItemsElemMatchFilter {
  and: [ProvidedAPIsItemsFilter]
  or: [ProvidedAPIsItemsFilter]
  nor: [ProvidedAPIsItemsFilter]
  not: ProvidedAPIsItemsFilter
}

input ProvidedAPIsItemsFilter {
  group: StringFilter
  kind: StringFilter
  plural: StringFilter
  version: StringFilter
  and: [ProvidedAPIsItemsFilter]
  or: [ProvidedAPIsItemsFilter]
  nor: [ProvidedAPIsItemsFilter]
  not: ProvidedAPIsItemsFilter
}

input RelatedImagesItemsElemMatchFilter {
  and: [RelatedImagesItemsFilter]
  or: [RelatedImagesItemsFilter]
  nor: [RelatedImagesItemsFilter]
  not: RelatedImagesItemsFilter
}

input RelatedImagesItemsFilter {
  digest: StringFilter
  image: StringFilter
  name: StringFilter
  and: [RelatedImagesItemsFilter]
  or: [RelatedImagesItemsFilter]
  nor: [RelatedImagesItemsFilter]
  not: RelatedImagesItemsFilter
}

type ContainerImageEdges {
  rpm_manifest: ContainerImageRPMManifestResponse
  vulnerabilities(
    page_size: Int = 50
    page: Int = 0
    filter: ContainerImageVulnerabilityFilter
    sort_by: [SortBy]
  ): ContainerImageVulnerabilityPaginatedResponse
  test_results(
    page_size: Int = 50
    page: Int = 0
    filter: CertProjectTestResultFilter
    sort_by: [SortBy]
  ): CertProjectTestResultPaginatedResponse
  artifacts(
    page_size: Int = 50
    page: Int = 0
    filter: CertProjectArtifactFilter
    sort_by: [SortBy]
  ): CertProjectArtifactPaginatedResponse
}

type ContainerImageRPMManifestResponse {
  data: ContainerImageRPMManifest
  error: ResponseError
}

# A containerImageRPMManifest contains all the RPM packages for a given containerImage
type ContainerImageRPMManifest {
  # The foreign key to containerImage._id.
  image_id: ObjectID
  object_type: String

  # Content manifest of this image. RPM content included in the image.
  rpms: [RpmsItems]

  # MongoDB unique _id
  _id: String

  # The date when the entry was created. Value is created automatically on creation.
  creation_date: DateTime

  # The date when the entry was last updated.
  last_update_date: DateTime
  edges: ContainerImageRPMManifestEdges
}

# RPM content of an image.
type RpmsItems {
  # RPM architecture.
  architecture: String

  # GPG key used to sign the RPM.
  gpg: String

  # RPM name.
  name: String

  # RPM name, version, release, and architecture.
  nvra: String

  # RPM release.
  release: String

  # Source RPM name.
  srpm_name: String

  # Source RPM NEVRA (name, epoch, version, release, architecture).
  srpm_nevra: String

  # RPM summary.
  summary: String

  # RPM version.
  version: String
}

type ContainerImageRPMManifestEdges {
  image: ContainerImageResponse
}

type ContainerImageVulnerabilityPaginatedResponse {
  data: [ContainerImageVulnerability]
  error: ResponseError
  page: Int
  page_size: Int
  total: Int
}

# Vulnerability present in the content that is installed in the image.
type ContainerImageVulnerability {
  # Advisory identifier.
  advisory_id: String

  # Advisory type (RHSA, RHBA, ...).
  advisory_type: String

  # ID of the CVE.
  cve_id: String

  #
  object_type: String

  # Array of package information applicable to this CVE.
  packages: [ContainerImageVulnerabilityPackage]

  # Date the CVE was made public.
  public_date: String

  # CVE severity.
  severity: String

  # MongoDB unique _id
  _id: String

  # The date when the entry was created. Value is created automatically on creation.
  creation_date: DateTime

  # The date when the entry was last updated.
  last_update_date: DateTime
  edges: ContainerImageVulnerabilityEdges
}

# Package information applicable to this CVE.
type ContainerImageVulnerabilityPackage {
  # The next image that fixes this vulnerability.
  fixed_by_image: ContainerImageVulnerabilityFixedByImage

  # The source RPM that contains the fix.
  srpm_nevra: String

  # RPMs, identified by their RPM NVRA, that are present in the image and vulnerable.
  rpm_nvra: [String]
}

# The image that fixes the vulnerability.
type ContainerImageVulnerabilityFixedByImage {
  # The _id of the image that fixes the vulnerability.
  id: String

  # The data is denormalized to not require lookup in UI.
  #
  # RISK: could get out of sync with real refrenced data.
  repositories: [ContainerImageVulnerabilityFixedByRepository]
  edges: ContainerImageVulnerabilityFixedByImageEdges
}

#
type ContainerImageVulnerabilityFixedByRepository {
  #
  registry: String

  #
  repository: String

  #
  tags: [ContainerImageVulnerabilityTag]
}

#
type ContainerImageVulnerabilityTag {
  #
  name: String
}

type ContainerImageVulnerabilityFixedByImageEdges {
  image: ContainerImageResponse
}

type ContainerImageVulnerabilityEdges {
  advisory: RedHatContainerAdvisoryResponse
}

input ContainerImageVulnerabilityFilter {
  advisory_id: StringFilter
  advisory_type: StringFilter
  cve_id: StringFilter
  object_type: StringFilter
  packages_size: ListSizeFilter
  packages_index: StringIndexFilter
  packages_elemMatch: ContainerImageVulnerabilityPackageElemMatchFilter
  packages: ContainerImageVulnerabilityPackageFilter
  public_date: StringFilter
  severity: StringFilter
  _id: StringFilter
  creation_date: DateTimeFilter
  last_update_date: DateTimeFilter
  and: [ContainerImageVulnerabilityFilter]
  or: [ContainerImageVulnerabilityFilter]
  nor: [ContainerImageVulnerabilityFilter]
  not: ContainerImageVulnerabilityFilter
}

input ContainerImageVulnerabilityPackageElemMatchFilter {
  and: [ContainerImageVulnerabilityPackageFilter]
  or: [ContainerImageVulnerabilityPackageFilter]
  nor: [ContainerImageVulnerabilityPackageFilter]
  not: ContainerImageVulnerabilityPackageFilter
}

input ContainerImageVulnerabilityPackageFilter {
  fixed_by_image: ContainerImageVulnerabilityFixedByImageFilter
  srpm_nevra: StringFilter
  rpm_nvra_size: ListSizeFilter
  rpm_nvra_index: StringIndexFilter
  rpm_nvra: StringFilter
  and: [ContainerImageVulnerabilityPackageFilter]
  or: [ContainerImageVulnerabilityPackageFilter]
  nor: [ContainerImageVulnerabilityPackageFilter]
  not: ContainerImageVulnerabilityPackageFilter
}

input ContainerImageVulnerabilityFixedByImageFilter {
  id: StringFilter
  repositories_size: ListSizeFilter
  repositories_index: StringIndexFilter
  repositories_elemMatch: ContainerImageVulnerabilityFixedByRepositoryElemMatchFilter
  repositories: ContainerImageVulnerabilityFixedByRepositoryFilter
  and: [ContainerImageVulnerabilityFixedByImageFilter]
  or: [ContainerImageVulnerabilityFixedByImageFilter]
  nor: [ContainerImageVulnerabilityFixedByImageFilter]
  not: ContainerImageVulnerabilityFixedByImageFilter
}

input ContainerImageVulnerabilityFixedByRepositoryElemMatchFilter {
  and: [ContainerImageVulnerabilityFixedByRepositoryFilter]
  or: [ContainerImageVulnerabilityFixedByRepositoryFilter]
  nor: [ContainerImageVulnerabilityFixedByRepositoryFilter]
  not: ContainerImageVulnerabilityFixedByRepositoryFilter
}

input ContainerImageVulnerabilityFixedByRepositoryFilter {
  registry: StringFilter
  repository: StringFilter
  tags_size: ListSizeFilter
  tags_index: StringIndexFilter
  tags_elemMatch: ContainerImageVulnerabilityTagElemMatchFilter
  tags: ContainerImageVulnerabilityTagFilter
  and: [ContainerImageVulnerabilityFixedByRepositoryFilter]
  or: [ContainerImageVulnerabilityFixedByRepositoryFilter]
  nor: [ContainerImageVulnerabilityFixedByRepositoryFilter]
  not: ContainerImageVulnerabilityFixedByRepositoryFilter
}

input ContainerImageVulnerabilityTagElemMatchFilter {
  and: [ContainerImageVulnerabilityTagFilter]
  or: [ContainerImageVulnerabilityTagFilter]
  nor: [ContainerImageVulnerabilityTagFilter]
  not: ContainerImageVulnerabilityTagFilter
}

input ContainerImageVulnerabilityTagFilter {
  name: StringFilter
  and: [ContainerImageVulnerabilityTagFilter]
  or: [ContainerImageVulnerabilityTagFilter]
  nor: [ContainerImageVulnerabilityTagFilter]
  not: ContainerImageVulnerabilityTagFilter
}

type CertProjectTestResultPaginatedResponse {
  data: [CertProjectTestResult]
  error: ResponseError
  page: Int
  page_size: Int
  total: Int
}

# Contain certification test results of related certProject
type CertProjectTestResult {
  # Unique identifier for the certification project.
  cert_project: ObjectID

  # Hashed content for the certification project.
  certification_hash: String

  # Image associated with the test result.
  image: String

  # Operator package name associated with the test result.
  operator_package_name: String

  # Red Hat Org ID / account_id from Red Hat SSO. Also corresponds to company_org_id in Red Hat Connect.
  org_id: Int

  # Whether or not the test has passed overall.
  passed: Boolean

  # Identifier of container image collection.
  image_id: ObjectID

  # The test results stored in lists based on result status.
  results: Results

  # The test library of the test result.
  test_library: TestLibrary

  # Version associated with the content tested.
  version: String

  # Pull request of certification test results
  pull_request: PullRequest

  # A platform where tests were executed.
  tested_on: TestedOn

  # MongoDB unique _id
  _id: String

  # The date when the entry was created. Value is created automatically on creation.
  creation_date: DateTime

  # The date when the entry was last updated.
  last_update_date: DateTime
  edges: CertProjectTestResultEdges
}

# The test results stored in lists based on result status.
type Results {
  # Test results of cert project certification
  failed: [TestResults]

  # Test results of cert project certification
  errors: [TestResults]

  # Test results of cert project certification
  passed: [TestResults]
}

# The cert project pipeline test result.
type TestResults {
  check_url: URI
  description: String
  elapsed_time: Float
  help: String
  knowledgebase_url: URI
  name: String
  suggestion: String
}

# The test library of the test result.
type TestLibrary {
  commit: String
  name: String
  version: String
}

# Pull request of certification test results.
type PullRequest {
  # Pull request URL
  url: String

  # Pull request identifier
  id: Int

  # Pull request status
  status: String
}

# A platform where tests were executed.
type TestedOn {
  name: String
  version: String
}

type CertProjectTestResultEdges {
  cert_project: CertificationProjectResponse
  container_image: ContainerImageResponse
}

input CertProjectTestResultFilter {
  cert_project: StringFilter
  certification_hash: StringFilter
  image: StringFilter
  operator_package_name: StringFilter
  org_id: IntFilter
  passed: BooleanFilter
  image_id: StringFilter
  results: ResultsFilter
  test_library: TestLibraryFilter
  version: StringFilter
  pull_request: PullRequestFilter
  tested_on: TestedOnFilter
  _id: StringFilter
  creation_date: DateTimeFilter
  last_update_date: DateTimeFilter
  and: [CertProjectTestResultFilter]
  or: [CertProjectTestResultFilter]
  nor: [CertProjectTestResultFilter]
  not: CertProjectTestResultFilter
}

input ResultsFilter {
  failed_size: ListSizeFilter
  failed_index: StringIndexFilter
  failed_elemMatch: TestResultsElemMatchFilter
  failed: TestResultsFilter
  errors_size: ListSizeFilter
  errors_index: StringIndexFilter
  errors_elemMatch: TestResultsElemMatchFilter
  errors: TestResultsFilter
  passed_size: ListSizeFilter
  passed_index: StringIndexFilter
  passed_elemMatch: TestResultsElemMatchFilter
  passed: TestResultsFilter
  and: [ResultsFilter]
  or: [ResultsFilter]
  nor: [ResultsFilter]
  not: ResultsFilter
}

input TestResultsElemMatchFilter {
  and: [TestResultsFilter]
  or: [TestResultsFilter]
  nor: [TestResultsFilter]
  not: TestResultsFilter
}

input TestResultsFilter {
  check_url: URIFilter
  description: StringFilter
  elapsed_time: FloatFilter
  help: StringFilter
  knowledgebase_url: URIFilter
  name: StringFilter
  suggestion: StringFilter
  and: [TestResultsFilter]
  or: [TestResultsFilter]
  nor: [TestResultsFilter]
  not: TestResultsFilter
}

input FloatFilter {
  eq: Float
  ne: Float
  le: Float
  lt: Float
  ge: Float
  gt: Float
  in: [Float]
}

input TestLibraryFilter {
  commit: StringFilter
  name: StringFilter
  version: StringFilter
  and: [TestLibraryFilter]
  or: [TestLibraryFilter]
  nor: [TestLibraryFilter]
  not: TestLibraryFilter
}

input PullRequestFilter {
  url: StringFilter
  id: IntFilter
  status: StringFilter
  and: [PullRequestFilter]
  or: [PullRequestFilter]
  nor: [PullRequestFilter]
  not: PullRequestFilter
}

input TestedOnFilter {
  name: StringFilter
  version: StringFilter
  and: [TestedOnFilter]
  or: [TestedOnFilter]
  nor: [TestedOnFilter]
  not: TestedOnFilter
}

type CertProjectArtifactPaginatedResponse {
  data: [CertProjectArtifact]
  error: ResponseError
  page: Int
  page_size: Int
  total: Int
}

# Contains the logs from the preflight test results from operator pipelines project
type CertProjectArtifact {
  # Cert Project Identifier.
  cert_project: ObjectID

  # Red Hat Org ID / account_id from Red Hat SSO. Also corresponds to company_org_id in Red Hat Connect.
  org_id: Int

  # Content version.
  version: String

  # The operator package name of the cert project artifact.
  operator_package_name: String

  # The cert project hashed content.
  certification_hash: String

  # Identifier of container image collection.
  image_id: ObjectID

  # Base64 encoded the cert project artifact content.
  content: String

  # The content type associated with the content type.
  content_type: String

  # The file name associated with the content test results.
  filename: String

  # File size in bytes.
  file_size: Int

  # MongoDB unique _id
  _id: String

  # The date when the entry was created. Value is created automatically on creation.
  creation_date: DateTime

  # The date when the entry was last updated.
  last_update_date: DateTime
  edges: CertProjectArtifactEdges
}

type CertProjectArtifactEdges {
  cert_project: CertificationProjectResponse
  container_image: ContainerImageResponse
}

input CertProjectArtifactFilter {
  cert_project: StringFilter
  org_id: IntFilter
  version: StringFilter
  operator_package_name: StringFilter
  certification_hash: StringFilter
  image_id: StringFilter
  content: StringFilter
  content_type: StringFilter
  filename: StringFilter
  file_size: IntFilter
  _id: StringFilter
  creation_date: DateTimeFilter
  last_update_date: DateTimeFilter
  and: [CertProjectArtifactFilter]
  or: [CertProjectArtifactFilter]
  nor: [CertProjectArtifactFilter]
  not: CertProjectArtifactFilter
}

type ContainerFilePaginatedResponse {
  data: [ContainerFile]
  error: ResponseError
  page: Int
  page_size: Int
  total: Int
}

#
type ContainerFile {
  filename: String
  content: String
  key: String
}

input ContainerFileFilter {
  filename: StringFilter
  content: StringFilter
  key: StringFilter
  and: [ContainerFileFilter]
  or: [ContainerFileFilter]
  nor: [ContainerFileFilter]
  not: ContainerFileFilter
}

type ContainerTagHistoryPaginatedResponse {
  data: [ContainerTagHistory]
  error: ResponseError
  page: Int
  page_size: Int
  total: Int
}

input ContainerTagHistoryFilter {
  object_type: StringFilter
  registry: StringFilter
  repository: StringFilter
  tag: StringFilter
  tag_type: StringFilter
  history_size: ListSizeFilter
  history_index: StringIndexFilter
  history_elemMatch: HistoryElemMatchFilter
  history: HistoryFilter
  _id: StringFilter
  creation_date: DateTimeFilter
  last_update_date: DateTimeFilter
  and: [ContainerTagHistoryFilter]
  or: [ContainerTagHistoryFilter]
  nor: [ContainerTagHistoryFilter]
  not: ContainerTagHistoryFilter
}

input HistoryElemMatchFilter {
  and: [HistoryFilter]
  or: [HistoryFilter]
  nor: [HistoryFilter]
  not: HistoryFilter
}

input HistoryFilter {
  brew_build: StringFilter
  end_date: DateTimeFilter
  start_date: DateTimeFilter
  and: [HistoryFilter]
  or: [HistoryFilter]
  nor: [HistoryFilter]
  not: HistoryFilter
}

type CertProjectBuildRequestResponse {
  data: CertProjectBuildRequest
  error: ResponseError
}

type CertProjectScanRequestResponse {
  data: CertProjectScanRequest
  error: ResponseError
}

type CertProjectTagRequestResponse {
  data: CertProjectTagRequest
  error: ResponseError
}

type CertProjectTestResultResponse {
  data: CertProjectTestResult
  error: ResponseError
}

type CertProjectArtifactResponse {
  data: CertProjectArtifact
  error: ResponseError
}

type OperatorPackageResponse {
  data: OperatorPackage
  error: ResponseError
}

#
type OperatorPackage {
  #
  association: String

  #
  package_name: String

  #
  source: String

  # MongoDB unique _id
  _id: ObjectID

  # The date when the entry was created. Value is created automatically on creation.
  creation_date: DateTime

  # The date when the entry was last updated.
  last_update_date: DateTime
}

type OperatorPackagePaginatedResponse {
  data: [OperatorPackage]
  error: ResponseError
  page: Int
  page_size: Int
  total: Int
}

input OperatorPackageFilter {
  association: StringFilter
  package_name: StringFilter
  source: StringFilter
  _id: StringFilter
  creation_date: DateTimeFilter
  last_update_date: DateTimeFilter
  and: [OperatorPackageFilter]
  or: [OperatorPackageFilter]
  nor: [OperatorPackageFilter]
  not: OperatorPackageFilter
}

type OperatorIndexResponse {
  data: OperatorIndex
  error: ResponseError
}

#
type OperatorIndex {
  # OCP version, e.g. 4.5.
  ocp_version: SemVer

  # Organization, as understood by iib, e.g. redhat-marketplace.
  organization: String

  # The docker path used to pull this index container, e.g. quay.io/foo/bar:v4.5.
  path: String

  # The date till the index image is valid
  end_of_life: DateTime

  # MongoDB unique _id
  _id: String

  # The date when the entry was created. Value is created automatically on creation.
  creation_date: DateTime

  # The date when the entry was last updated.
  last_update_date: DateTime
}

type OperatorIndexPaginatedResponse {
  data: [OperatorIndex]
  error: ResponseError
  page: Int
  page_size: Int
  total: Int
}

input OperatorIndexFilter {
  ocp_version: StringFilter
  organization: StringFilter
  path: StringFilter
  end_of_life: DateTimeFilter
  _id: StringFilter
  creation_date: DateTimeFilter
  last_update_date: DateTimeFilter
  and: [OperatorIndexFilter]
  or: [OperatorIndexFilter]
  nor: [OperatorIndexFilter]
  not: OperatorIndexFilter
}

type RedHatContainerAdvisoryPaginatedResponse {
  data: [RedHatContainerAdvisory]
  error: ResponseError
  page: Int
  page_size: Int
  total: Int
}

input RedHatContainerAdvisoryFilter {
  _id: StringFilter
  content_type: StringFilter
  description: StringFilter
  object_type: StringFilter
  severity: StringFilter
  ship_date: DateTimeFilter
  solution: StringFilter
  synopsis: StringFilter
  topic: StringFilter
  type: StringFilter
  cves_size: ListSizeFilter
  cves_index: StringIndexFilter
  cves_elemMatch: CVEElemMatchFilter
  cves: CVEFilter
  issues_size: ListSizeFilter
  issues_index: StringIndexFilter
  issues_elemMatch: IssueElemMatchFilter
  issues: IssueFilter
  creation_date: DateTimeFilter
  last_update_date: DateTimeFilter
  and: [RedHatContainerAdvisoryFilter]
  or: [RedHatContainerAdvisoryFilter]
  nor: [RedHatContainerAdvisoryFilter]
  not: RedHatContainerAdvisoryFilter
}

input CVEElemMatchFilter {
  and: [CVEFilter]
  or: [CVEFilter]
  nor: [CVEFilter]
  not: CVEFilter
}

input CVEFilter {
  id: StringFilter
  url: StringFilter
  and: [CVEFilter]
  or: [CVEFilter]
  nor: [CVEFilter]
  not: CVEFilter
}

input IssueElemMatchFilter {
  and: [IssueFilter]
  or: [IssueFilter]
  nor: [IssueFilter]
  not: IssueFilter
}

input IssueFilter {
  id: StringFilter
  issue_tracker: StringFilter
  url: StringFilter
  and: [IssueFilter]
  or: [IssueFilter]
  nor: [IssueFilter]
  not: IssueFilter
}

type BinarySignatureResponse {
  data: BinarySignature
  error: ResponseError
}

# Object for GraphQL response
type BinarySignature {
  binary_signature: String
}

# Combination of field name and the sorting direction used to sort the responses. If multiple pairs are set, they go from the most important to the least important.
input SortByMembersInput {
  # Name of the field that should be used to sort results. If the field is nested, use dot notation.
  field: String

  # If the fields should be in ascending or descending order.
  order: String
}

type ContainerImageVulnerabilityResponse {
  data: ContainerImageVulnerability
  error: ResponseError
}

type ProductListingResponse {
  data: ProductListing
  error: ResponseError
}

type AnalyticsPageViewsResponse {
  data: AnalyticsPageViews
  error: ResponseError
}

# Page views statistics.
type AnalyticsPageViews {
  # Page view statistics by date.
  by_date: [AnalyticsPageViewsByDate]

  # Total number of page views.
  total_pageviews: Int
}

# Page view statistics by date.
type AnalyticsPageViewsByDate {
  # Date of the page view.
  activity_date: DateTime

  # Number of page views.
  pageviews: Int
}

type AnalyticsPullCountResponse {
  data: AnalyticsPullCount
  error: ResponseError
}

# Pull count statistics.
type AnalyticsPullCount {
  # Pull count statistics by costumer.
  by_customers: [AnalyticsPullCountByCostumer]

  # Pull count statistics by tags.
  by_tags: [AnalyticsPullCountByTag]

  # Total number of distinct customers.
  total_customers: Int

  # Total number of distinct countries customers are from.
  total_countries: Int

  # Total number of image pulls performed.
  total_pulls: Int
}

# Pull count statistics by costumer.
type AnalyticsPullCountByCostumer {
  # Country that the customer is from.
  country: String

  # Name of the customer.
  customer_name: String

  # Date of the image pull.
  download_date: DateTime

  # Number of image pulls associated with the customer.
  pull_count: Int
}

# Pull count statistics by tags.
type AnalyticsPullCountByTag {
  # Date of the image pull.
  download_date: DateTime

  # Tags of the image pulled.
  image_tags: [String]

  # Number of image pulls associated with their tags.
  pull_count: Int
}

type ApiKeyListResponse {
  data: [ApiKey]
  error: ResponseError
}

# API key stored in Loki.
type ApiKey {
  id: Int
  description: String
  company_id: Int
  created: DateTime
  last_used: DateTime
  created_by: String

  # Red Hat Org ID / account_id from Red Hat SSO. Also corresponds to company_org_id in Red Hat Connect.
  org_id: Int
}

type ContainerGradesListResponse {
  data: [ContainerGrades]
  error: ResponseError
}

# Grades are describing security risk with containers that Red Hat provides through the Red Hat Ecosystem Catalog.
type ContainerGrades {
  # Image architecture.
  architecture: String

  # Current image security grade based on current date and grading schedule.
  current_grade: String

  # the grade based on applicable updates and time provided by PST CVE engine.
  freshness_grades: [FreshnessGrade]

  # Unique identifier for image.
  image_id: String

  # A date when current grade drops.
  next_drop_date: DateTime

  # Name of floating tag associated with the image.
  tag: String
  edges: ContainerGradesEdges
}

type ContainerGradesEdges {
  image: ContainerImageResponse
}

type Mutation {
  # Create image.
  create_image(input: ContainerImageInput): ContainerImageResponse

  # Update/Patch image.
  update_image(id: String, input: ContainerImageInput): ContainerImageResponse

  # Replace image.
  replace_image(id: String, input: ContainerImageInput): ContainerImageResponse

  # Replace container image by manifest digest.
  put_image_by_manifest_digest_registry_and_repository(
    repository: String
    registry: String
    manifest_digest: String
    input: ContainerImageInput
  ): ContainerImageResponse

  # Update container image by manifest digest.
  patch_image_by_manifest_digest_registry_and_repository(
    repository: String
    registry: String
    manifest_digest: String
    input: ContainerImageInput
  ): ContainerImageResponse

  # Update/Patch a RPM manifest by ID.
  update_image_rpm_manifest(
    id: String
    input: ContainerImageRPMManifestInput
  ): ContainerImageRPMManifestResponse

  # Replace a RPM Manifest by ID.
  replace_image_rpm_manifest(
    id: String
    input: ContainerImageRPMManifestInput
  ): ContainerImageRPMManifestResponse

  # Create a new RPM manifest for an image.
  create_image_rpm_manifest(
    id: String
    input: ContainerImageRPMManifestInput
  ): ContainerImageRPMManifestResponse

  # Create a certification project build request
  create_certification_project_build_request(
    id: String
    input: CertProjectBuildRequestInput
  ): CertProjectBuildRequestResponse

  # Create a certification project scan request
  create_certification_project_scan_request(
    id: String
    input: CertProjectScanRequestInput
  ): CertProjectScanRequestResponse

  # Create a certification project tag request
  create_certification_project_tag_request(
    id: String
    input: CertProjectTagRequestInput
  ): CertProjectTagRequestResponse

  # Partially update a vendor.
  update_vendor(
    id: String
    input: ContainerVendorInput
  ): ContainerVendorResponse

  # Create a certification project.
  create_certification_project(
    input: CertificationProjectInput
  ): CertificationProjectResponse

  # Partially update a certification project.
  update_certification_project(
    id: String
    input: CertificationProjectInput
  ): CertificationProjectResponse

  # Update a certification project.
  replace_certification_project(
    id: String
    input: CertificationProjectInput
  ): CertificationProjectResponse

  # Replace product listing.
  replace_product_listing(
    id: String
    input: ProductListingInput
  ): ProductListingResponse

  # Update product listing.
  update_product_listing(
    id: String
    input: ProductListingInput
  ): ProductListingResponse

  # Create product listing.
  create_product_listing(input: ProductListingInput): ProductListingResponse

  # Create an API key.
  create_api_key(input: ApiKeyInput): ApiKeyResponse

  # Delete API key.
  delete_api_key(key_id: Int): ApiKeyResponse

  # Create a certification project test result
  create_certification_project_test_result(
    id: String
    input: CertProjectTestResultInput
  ): CertProjectTestResultResponse

  # Update/Patch certification project test result
  update_certification_project_test_result(
    id: String
    input: CertProjectTestResultInput
  ): CertProjectTestResultResponse

  # Create a certification project artifact
  create_certification_project_artifact(
    id: String
    input: CertProjectTestResultInput
  ): CertProjectArtifactResponse
}

# Metadata about images contained in RedHat and ISV repositories
input ContainerImageInput {
  # The field contains an architecture for which the container image was built for. Value is used to distinguish between the default x86-64 architecture and other architectures. If the value is not set, the image was built for the x86-64 architecture.
  architecture: String

  # Brew related metadata.
  brew: BrewInput

  # A list of all content sets (YUM repositories) from where an image RPM content is.
  content_sets: [String]

  # A mapping of applicable advisories to RPM NEVRA. This data is required for scoring.
  cpe_ids: [String]

  # A mapping of applicable advisories for the base_images from the Red Hat repositories.
  cpe_ids_rh_base_images: [String]

  # Docker Image Digest. For Docker 1.10+ this is also known as the 'manifest digest'.
  docker_image_digest: String

  # Docker Image ID. For Docker 1.10+ this is also known as the 'config digest'.
  docker_image_id: String

  # The grade based on applicable updates and time provided by PST CVE engine.
  freshness_grades: [FreshnessGradeInput]
  object_type: String

  # Data parsed from image metadata.
  # These fields are not computed from any other source.
  parsed_data: ParsedDataInput

  # Published repositories associated with the container image.
  repositories: [ContainerImageRepoInput]

  # Indication if the image was certified.
  certified: Boolean

  # Indicates that an image was removed. Only unpublished images can be removed.
  deleted: Boolean

  # Image manifest digest.
  # Be careful, as this value is not unique among container image entries, as one image can be references several times.
  image_id: String

  # ID of the project in for ISV repositories. The ID can be also used to connect vendor to the image.
  isv_pid: String

  # The total size of the sum of all layers for each image in bytes. This is computed externally and may not match what is reported by the image metadata (see parsed_data.size).
  sum_layer_size_bytes: Int

  # Field for multiarch primary key
  top_layer_id: String

  # Hash (sha256) of the uncompressed top layer for this image (should be same value as - parsed_data.uncompressed_layer_sizes.0.layer_id)
  uncompressed_top_layer_id: String

  # Raw image configuration, such as output from docker inspect.
  raw_config: String

  # MongoDB unique _id
  _id: String

  # The date when the entry was created. Value is created automatically on creation.
  creation_date: DateTime

  # The date when the entry was last updated.
  last_update_date: DateTime
}

# Brew Build System related metadata.
input BrewInput {
  # Unique and immutable Brew build ID.
  build: String

  # Timestamp from Brew when the image has been succesfully built.
  completion_date: DateTime

  # Multi-Arch primary key.
  nvra: String

  # A package name in Brew.
  package: String
}

#
input CertificationInput {
  assessment: [AssessmentInput]
}

#
input AssessmentInput {
  # Assesment name.
  name: String

  # Indicates if the assessment is required for certification.
  required_for_certification: Boolean

  # Indicates if the assesment was passed, True means yes.
  value: Boolean
}

# Grade based on applicable updates and time provided by PST CVE engine.
input FreshnessGradeInput {
  # Date after which the grade is no longer valid. See start_date for when the grade was effective. If no value is set, the grade applies forever. This should happen only for a grade of A (no vulnerabilities) or grade F.
  end_date: DateTime

  # The grade.
  grade: String

  # Date when the grade was added by the vulnerability engine.
  creation_date: DateTime

  # Date from which the grade is in effect. The grade is effective until the end_date, if end_date is set.
  start_date: DateTime
}

#
input ParsedDataInput {
  architecture: String
  author: String
  command: String
  comment: String
  container: String

  # The 'created' date reported by image metadata. Stored as String because we do not have control on that format.
  created: String
  docker_image_digest: String
  docker_image_id: String

  # Version of docker reported by 'docker inspect' for this image.
  docker_version: String
  env_variables: [String]
  image_id: String
  labels: [LabelInput]

  # Layer digests from the image.
  layers: [String]
  os: String
  ports: String

  # Repositories defined within an image as reported by yum command.
  repos: [ParsedDataRepoInput]

  # Size of this image as reported by image metadata.
  size: Int

  # Information about uncompressed layer sizes.
  uncompressed_layer_sizes: [UncompressedLayerSizeInput]

  # Uncompressed images size in bytes (sum of uncompressed layers size).
  uncompressed_size_bytes: Int

  # The user on the images.
  user: String

  # Virtual size of this image as reported by image metadata.
  virtual_size: Int
}

# Image label.
input LabelInput {
  # The name of the label
  name: String

  # Value of the label.
  value: String
}

#
input ParsedDataRepoInput {
  baseurl: String
  expire: String
  filename: String
  id: String
  name: String
  pkgs: String
  size: String
  updated: String
}

#
input UncompressedLayerSizeInput {
  # The SHA256 layer ID.
  layer_id: String

  # The uncompressed layer size in bytes.
  size_bytes: Int
}

#
input ContainerImageRepoInput {
  # Store information about image comparison.
  comparison: ContainerImageRepoComparisonInput

  # The _id's of the redHatContainerAdvisory that contains the content advisories.
  content_advisory_ids: [String]

  # The _id of the redHatContainerAdvisory that contains the image advisory.
  image_advisory_id: String

  # Available for multiarch images.
  manifest_list_digest: String

  # Available for single arch images.
  manifest_schema2_digest: String

  # Indicate if the image has been published to the container catalog.
  published: Boolean

  # Date the image was published to the container catalog.
  published_date: DateTime

  # When the image was pushed to this repository. For RH images this is picked from first found of advisory ship_date, brew completion_date, and finally repositories publish_date. For ISV images this TBD but is probably going to be only sourced from publish_date but could come from parsed_data.created.
  push_date: DateTime

  # Hostname of the registry where the repository can be accessed.
  registry: String

  # Repository name.
  repository: String

  # Image signing info.
  signatures: [SignatureInfoInput]

  # List of container tags assigned to this layer.
  tags: [ContainerImageRepoTagInput]
}

#
input ContainerImageRepoComparisonInput {
  # Mapping of a NVRA to multiple advisories IDs.
  advisory_rpm_mapping: [ContainerImageRepoComparisonMappingInput]

  # Reason why 'with_nvr' is or is not null.
  reason: String

  # Human readable reason.
  reason_text: String

  # List of rpms grouped by category (new, remove, upgrade, downgrade).
  rpms: ContainerImageRepoComparisonRPMsInput

  # NVR of image which this image was compared with.
  with_nvr: String
}

#
input ContainerImageRepoComparisonMappingInput {
  # Content advisory ID.
  advisory_ids: [String]

  # NVRA of the RPM related to advisories.
  nvra: String
}

#
input ContainerImageRepoComparisonRPMsInput {
  # List of NVRA which were downgraded in this image.
  downgrade: [String]

  # List of NVRA which were added to this image.
  new: [String]

  # List of NVRA which were removed in this image.
  remove: [String]

  # List of NVRA which were upgraded in this image.
  upgrade: [String]
}

#
input SignatureInfoInput {
  # The long 16-byte gpg key id.
  key_long_id: String

  # List of image tags that are signed with the given key.
  tags: [String]
}

#
input ContainerImageRepoTagInput {
  added_date: DateTime

  # Available when manifest_schema2_digest is not. All legacy images.
  manifest_schema1_digest: String

  # The name of the tag.
  name: String

  # Date this tag was removed from the image in this repo. If the tag is added back, add a new entry in 'tags' array.
  removed_date: DateTime
}

# A containerImageRPMManifest contains all the RPM packages for a given containerImage
input ContainerImageRPMManifestInput {
  # The foreign key to containerImage._id.
  image_id: ObjectID
  object_type: String

  # Content manifest of this image. RPM content included in the image.
  rpms: [RpmsItemsInput]

  # MongoDB unique _id
  _id: String

  # The date when the entry was created. Value is created automatically on creation.
  creation_date: DateTime

  # The date when the entry was last updated.
  last_update_date: DateTime
}

# RPM content of an image.
input RpmsItemsInput {
  # RPM architecture.
  architecture: String

  # GPG key used to sign the RPM.
  gpg: String

  # RPM name.
  name: String

  # RPM name, version, release, and architecture.
  nvra: String

  # RPM release.
  release: String

  # Source RPM name.
  srpm_name: String

  # Source RPM NEVRA (name, epoch, version, release, architecture).
  srpm_nevra: String

  # RPM summary.
  summary: String

  # RPM version.
  version: String
}

# Contain status and related metadata of a certProject build request.
input CertProjectBuildRequestInput {
  # Unique identifier for the certification project.
  cert_project: ObjectID

  # The request status
  status: String

  # The tag that the container image gets when build is done.
  tag: String

  # MongoDB unique _id
  _id: String

  # The date when the entry was created. Value is created automatically on creation.
  creation_date: DateTime

  # The date when the entry was last updated.
  last_update_date: DateTime

  # Red Hat Org ID / account_id from Red Hat SSO. Also corresponds to company_org_id in Red Hat Connect.
  org_id: Int

  # An explanatory message to a request status.
  status_message: String
}

# Contain status and related metadata of a certProject scan request.
input CertProjectScanRequestInput {
  # Unique identifier for the certification project.
  cert_project: ObjectID

  # URL pointing to the location of DCI logs.
  external_tests_link: URI

  # Image pull specification in repo@sha256:digest format.
  pull_spec: String

  # Unique identifier of an ISV certification scan
  scan_uuid: String

  # Container image tag associated with the scan request.
  tag: String

  # MongoDB unique _id
  _id: String

  # The date when the entry was created. Value is created automatically on creation.
  creation_date: DateTime

  # The date when the entry was last updated.
  last_update_date: DateTime

  # Red Hat Org ID / account_id from Red Hat SSO. Also corresponds to company_org_id in Red Hat Connect.
  org_id: Int

  # The request status
  status: String

  # An explanatory message to a request status.
  status_message: String
}

# Contain status and related metadata of a certProject tag request.
input CertProjectTagRequestInput {
  # Unique identifier for the certification project.
  cert_project: ObjectID

  # Container image id associated with the tag request.
  image_id: ObjectID

  # Operation performed during the tag request, e.g. publish
  operation: String

  # Container image tag associated with the tag request.
  tag: String

  # MongoDB unique _id
  _id: String

  # The date when the entry was created. Value is created automatically on creation.
  creation_date: DateTime

  # The date when the entry was last updated.
  last_update_date: DateTime

  # Red Hat Org ID / account_id from Red Hat SSO. Also corresponds to company_org_id in Red Hat Connect.
  org_id: Int

  # The request status
  status: String

  # An explanatory message to a request status.
  status_message: String
}

# Stores information about a Vendor
input ContainerVendorInput {
  # URL to the vendor's main website.
  company_url: URI

  # General contact information for the vendor, to be displayed on the vendor page on RHCC.
  contact: ContainerVendorContactInput
  description: String

  # Company node ID from Red Hat Connect.
  drupal_company_id: Int

  # The industry / vertical the vendor belongs to.
  industries: [String]
  label: String

  # A flag that determines if vendor label can be changed.
  label_locked: Boolean
  logo_url: URI
  name: String
  object_type: String

  # Indicate that the vendor has been published.
  published: Boolean
  registry_urls: [String]

  # RSS feed for vendor.
  rss_feed_url: URI

  # Token for outbound namespace for pulling published marketplace images.
  service_account_token: String
  social_media_links: [ContainerVendorSocialMediaLinksInput]

  # Red Hat Org ID / account_id from Red Hat SSO. Also corresponds to company_org_id in Red Hat Connect.
  org_id: Int

  # MongoDB unique _id
  _id: String

  # The date when the entry was created. Value is created automatically on creation.
  creation_date: DateTime

  # The date when the entry was last updated.
  last_update_date: DateTime
}

# Contact information
input ContainerVendorContactInput {
  # General contact email address.
  email: String

  # General contact phone number.
  phone: String
}

# Social media links.
input ContainerVendorSocialMediaLinksInput {
  # The name of the social media provider.
  name: String

  # The URL to the social media site for the vendor.
  url: URI
}

# Certification project information.
input CertificationProjectInput {
  # Certification Date.
  certification_date: DateTime

  # Certification Status.
  certification_status: String

  # Certification User.
  certification_user: Int

  # Contacts for certification project.
  contacts: [CertProjectContactsInput]
  container: CertProjectContainerInput

  # Configuration specific to Helm Chart projects.
  helm_chart: CertProjectHelmChartInput
  marketplace: CertProjectMarketplaceInput

  # The owner provided name of the certification project.
  name: String

  # Operator Distribution.
  operator_distribution: String

  # Red Hat Org ID / account_id from Red Hat SSO. Also corresponds to company_org_id in Red Hat Connect.
  org_id: Int

  # Unique identifier for the product listing.
  product_listings: [String]

  # Status of the certification project.
  project_status: String

  # Who published the certification project.
  published_by: String
  redhat: CertProjectRedhatInput
  self_certification: CertProjectSelfCertificationInput

  # Certification project type.
  type: String

  # MongoDB unique _id
  _id: String

  # The date when the entry was created. Value is created automatically on creation.
  creation_date: DateTime

  # The date when the entry was last updated.
  last_update_date: DateTime
}

# Contact info.
input CertProjectContactsInput {
  email_address: Email
  type: String
}

# Container related information.
input CertProjectContainerInput {
  # The application categories (types).
  application_categories: [String]

  # Once a container is certified it is automatically published. Auto-publish must be enabled in order to set up automatic rebuilds. Auto-publish must always be enabled when auto-rebuilding is enabled.
  auto_publish: Boolean

  # Auto rebuild enabled.
  auto_rebuild: Boolean

  # Distribution approval obtained.
  distribution_approval: Boolean

  # Distribution method.
  distribution_method: String

  # ID of the project in for ISV repositories.
  isv_pid: String

  # Kubernetes objects for operator registry projects. Value has to be a valid YAML.
  kube_objects: OpenPGPEncrypted

  # Docker config for operator registry projects. Value has to be a valid JSON.
  docker_config_json: OpenPGPEncrypted

  # OS Content Type.
  os_content_type: String

  # Passed RH Cert.
  passed_rhcert: Boolean

  # A container needs to run in a privileged state.
  privileged: Boolean

  # Hostname of the registry where the repository can be accessed.
  # Examples: registry.company.com assumes the default port, 443. registry.company.com:5000 repository path with optional port specified.
  # It is only applicable for projects with an 'external' distribution method.
  registry: String

  # Note: These instructions will be displayed in the Red Hat Container Catalog as is. Please modify the following template as it suits your needs.
  registry_override_instruct: String

  # Release category.
  release_category: String

  # Path to the container repository as found in the registry.
  #
  # Examples:
  # path/to/repository
  # repository
  #
  # This field can only be edited when there are no published containers.
  # It is only applicable for projects with an 'external' distribution method.
  repository: String

  # The repository description is displayed on the container
  # catalog repository overview page.
  repository_description: String

  # This should represent your product (or the component if your product consists of multiple containers)
  # and a major version. For example, you could use names like jboss-server7, or agent5.
  #
  # This value is only editable when there are no published containers in this project.
  # It is only applicable for projects that do not have the 'external' distribution method.
  repository_name: String

  # Service Account Secret.
  service_account_secret: String

  # Short description of the container.
  short_description: String

  # Supported Platforms.
  support_platforms: [String]

  # Container type.
  # Field is required, if project type is 'Container', and the field is immutable for Partners after creation.
  type: String

  # Filename other than the default Dockerfile or a path to a Dockerfile in a subdirectory.
  source_dockerfile: String

  # Force the build to ignore cached layers and rerun all steps of the Dockerfile.
  build_no_cache: Boolean

  # Override default location (root directory) for applications within a subdirectory.
  source_context_dir: String

  # Whether Red Hat will build your container.
  build_service: Boolean

  # The specific Git branch to checkout.
  source_ref: String

  # The URL to the source used for the build.
  # For example: 'https://github.com/openshift/ruby-hello-world
  source_uri: URI

  # Base64 encoded SSH private key in PEM format. Used to pull the source.
  source_ssh_private_key: Base64OpenPGPEncrypted

  # GitHub users authorized to submit a certification pull request.
  github_usernames: [String]
}

# Helm chart related information.
input CertProjectHelmChartInput {
  # How your Helm Chart is distributed.
  distribution_method: String

  # The Helm Chart name as it will appear in GitHub.
  chart_name: String

  # URL to the externally distributed Helm Chart repository. This is not used if the chart is distributed via Red Hat.
  repository: URI

  # Instructions for users to access an externally distributed Helm Chart.
  distribution_instructions: String

  # Base64 encoded PGP public key. Used to sign result submissions.
  public_pgp_key: String

  # URL to the user submitted github pull request for this project.
  github_pull_request: URI

  # Short description of the Helm Chart.
  short_description: String

  # Long description of the Helm Chart.
  long_description: String

  # The application categories (types).
  application_categories: [String]

  # GitHub users authorized to submit a certification pull request.
  github_usernames: [String]
}

# Drupal related information.
input CertProjectDrupalInput {
  # Company node ID from Red Hat Connect.
  company_id: Int

  # Relation ID for certification project.
  relation: Int

  # Zone for certification project.
  zone: String
}

# Marketplace related information.
input CertProjectMarketplaceInput {
  enablement_status: String
  enablement_url: URI
  listing_url: URI
  published: Boolean
}

# Red Hat projects related information.
input CertProjectRedhatInput {
  # Red Hat Product ID.
  product_id: Int

  # Red Hat product name.
  product_name: String

  # Red Hat Product Version.
  product_version: String

  # Red Hat Product Version.
  product_version_id: Int
}

# Red Hat projects related information.
input CertProjectSelfCertificationInput {
  # Application Profiler.
  app_profiler: Boolean

  # Application Runs on App Type.
  app_runs_on_app_type: Boolean

  # Whether the Self Certification Evidence URL requires a customer login.
  auth_login: Boolean

  # Self Certification Evidence URL.
  certification_url: URI

  # Can Commercially Support on App Type.
  comm_support_on_app_type: Boolean

  # Self Certification Requested.
  requested: Boolean

  # TsaNET Member.
  tsanet_member: Boolean
}

# Product listings define a marketing page in the Ecosystem Catalog. It allows you to group repos and showcase what they accomplish together as an application. In the case of operators, your CSV file populates OperatorHub, which can only be viewed in cluster through OpenShift. Your product listing is publicly visible in the Ecosystem Catalog so anyone can know that it is offered.
input ProductListingInput {
  category: String

  # List of unique identifiers for the certification project.
  cert_projects: [String]

  # This field is required when the product listing is published.
  contacts: [ContactsItemsInput]

  # This field is required when the product listing is published.
  descriptions: DescriptionsInput

  # Company node ID from Red Hat Connect. Read only.
  drupal_company_id: Int

  # This field is required when the product listing is published.
  faqs: [FAQSItemsInput]

  # This field is required when the product listing is published.
  features: [FeaturesItemsInput]

  # This field is required when the product listing is published.
  functional_categories: [String]
  legal: LegalInput

  # This field is required when the product listing is published.
  linked_resources: [LinkedResourcesItemsInput]
  logo: LogoInput
  name: String
  published: Boolean

  # Flag determining if product listing is considered to be deleted. Product listing can be deleted only if it is not published. Value is set to False by default.
  deleted: Boolean
  quick_start_configuration: QuickStartConfigurationInput

  # List of unique identifiers for the repository.
  repositories: [String]

  # This field is required when the product listing is published.
  search_aliases: [SearchAliasesItemsInput]
  support: SupportInput
  type: String
  vendor_label: String
  operator_bundles: [OperatorBundlesItemsInput]

  # Red Hat Org ID / account_id from Red Hat SSO. Also corresponds to company_org_id in Red Hat Connect.
  org_id: Int

  # MongoDB unique _id
  _id: String

  # The date when the entry was created. Value is created automatically on creation.
  creation_date: DateTime

  # The date when the entry was last updated.
  last_update_date: DateTime
}

#
input BadgesItemsInput {
  badge: String
  project_id: ObjectID
}

#
input ContactsItemsInput {
  email_address: Email
  type: String
}

# This field is required when the product listing is published.
input DescriptionsInput {
  long: String
  short: String
}

# This field is required when the product listing is published.
input FAQSItemsInput {
  answer: String
  question: String
}

# This field is required when the product listing is published.
input FeaturesItemsInput {
  description: String
  title: String
}

#
input LegalInput {
  description: String
  license_agreement_url: URI
  privacy_policy_url: URI
}

#
input LinkedResourcesItemsInput {
  category: String
  description: String
  thumbnail_url: URI
  title: String
  type: String
  url: URI
}

#
input LogoInput {
  url: URI
}

#
input MarketplaceInput {
  enablement_status: String
  enablement_url: URI
  listing_url: URI
  published: Boolean
}

#
input QuickStartConfigurationInput {
  instructions: String
}

#
input SearchAliasesItemsInput {
  key: String
  value: String
}

# This field is required when the product listing is published.
input SupportInput {
  description: String
  email_address: Email
  phone_number: String
  url: URI
}

#
input OperatorBundlesItemsInput {
  # Bundle unique identifier
  _id: ObjectID

  # Bundle package name
  package: String
  capabilities: [String]
}

#
type ApiKeyResponse {
  api_key: String
  key_data: ApiKey
}

# API key stored in Loki.
input ApiKeyInput {
  id: Int
  description: String
  company_id: Int
  created: DateTime
  last_used: DateTime
  created_by: String

  # Red Hat Org ID / account_id from Red Hat SSO. Also corresponds to company_org_id in Red Hat Connect.
  org_id: Int
}

# Contain certification test results of related certProject
input CertProjectTestResultInput {
  # Unique identifier for the certification project.
  cert_project: ObjectID

  # Hashed content for the certification project.
  certification_hash: String

  # Image associated with the test result.
  image: String

  # Operator package name associated with the test result.
  operator_package_name: String

  # Red Hat Org ID / account_id from Red Hat SSO. Also corresponds to company_org_id in Red Hat Connect.
  org_id: Int

  # Whether or not the test has passed overall.
  passed: Boolean

  # Identifier of container image collection.
  image_id: ObjectID

  # The test results stored in lists based on result status.
  results: ResultsInput

  # The test library of the test result.
  test_library: TestLibraryInput

  # Version associated with the content tested.
  version: String

  # Pull request of certification test results
  pull_request: PullRequestInput

  # A platform where tests were executed.
  tested_on: TestedOnInput

  # MongoDB unique _id
  _id: String

  # The date when the entry was created. Value is created automatically on creation.
  creation_date: DateTime

  # The date when the entry was last updated.
  last_update_date: DateTime
}

# The test results stored in lists based on result status.
input ResultsInput {
  # Test results of cert project certification
  failed: [TestResultsInput]

  # Test results of cert project certification
  errors: [TestResultsInput]

  # Test results of cert project certification
  passed: [TestResultsInput]
}

# The cert project pipeline test result.
input TestResultsInput {
  check_url: URI
  description: String
  elapsed_time: Float
  help: String
  knowledgebase_url: URI
  name: String
  suggestion: String
}

# The test library of the test result.
input TestLibraryInput {
  commit: String
  name: String
  version: String
}

# Pull request of certification test results.
input PullRequestInput {
  # Pull request URL
  url: String

  # Pull request identifier
  id: Int

  # Pull request status
  status: String
}

# A platform where tests were executed.
input TestedOnInput {
  name: String
  version: String
}
`
