---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "backstage_group Data Source - terraform-provider-backstage"
subcategory: ""
description: |-
  Use this data source to get a specific Group entity https://backstage.io/docs/features/software-catalog/descriptor-format#kind-group from Backstage Software Catalog.
---

# backstage_group (Data Source)

Use this data source to get a specific [Group entity](https://backstage.io/docs/features/software-catalog/descriptor-format#kind-group) from Backstage Software Catalog.

## Example Usage

```terraform
# Retrieves specific group data:
data "backstage_group" "example" {
  # Required name of the group:
  name = "example-group"
  # If not provided, namespace defaults to "default" or the the one set in the provider:
  namespace = "example-namespace"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Name of the entity.

### Optional

- `namespace` (String) Namespace that the entity belongs to.

### Read-Only

- `api_version` (String) Version of specification format for this particular entity that this is written against.
- `id` (String) A globally unique ID for the entity. This field can not be set by the user at creation time, and the server will reject an attempt to do so. The field will be populated in read operations.
- `kind` (String) The high level entity type being described.
- `metadata` (Attributes) Metadata fields common to all versions/kinds of entity. (see [below for nested schema](#nestedatt--metadata))
- `relations` (Attributes List) Relations that this entity has with other entities (see [below for nested schema](#nestedatt--relations))
- `spec` (Attributes) The specification data describing the entity itself. (see [below for nested schema](#nestedatt--spec))

<a id="nestedatt--metadata"></a>
### Nested Schema for `metadata`

Read-Only:

- `annotations` (Map of String) Key/Value pairs of non-identifying auxiliary information attached to entity.
- `description` (String) A short (typically relatively few words) description of the entity.
- `etag` (String) An opaque string that changes for each update operation to any part of the entity, including metadata. This field can not be set by the user at creation time, and the server will reject an attempt to do so. The field will be populated in read operations.The field can (optionally) be specified when performing update or delete operations, and the server will then reject the operation if it does not match the current stored value.
- `labels` (Map of String) Key/Value pairs of identifying information attached to the entity.
- `links` (Attributes List) A list of external hyperlinks related to the entity. Links can provide additional contextual information that may be located outside of Backstage itself. For example, an admin dashboard or external CMS page. (see [below for nested schema](#nestedatt--metadata--links))
- `name` (String) Name of the entity.
- `namespace` (String) Namespace that the entity belongs to.
- `tags` (List of String) A list of single-valued strings, to for example classify catalog entities in various ways.
- `title` (String) A display name of the entity, to be presented in user interfaces instead of the name property, when available.
- `uid` (String) A globally unique ID for the entity. This field can not be set by the user at creation time, and the server will reject an attempt to do so. The field will be populated in read operations.

<a id="nestedatt--metadata--links"></a>
### Nested Schema for `metadata.links`

Read-Only:

- `icon` (String) A key representing a visual icon to be displayed in the UI.
- `title` (String) A user-friendly display name for the link.
- `type` (String) An optional value to categorize links into specific groups.
- `url` (String) URL in a standard uri format.



<a id="nestedatt--relations"></a>
### Nested Schema for `relations`

Read-Only:

- `target` (Attributes) The entity of the target of this relation. (see [below for nested schema](#nestedatt--relations--target))
- `target_ref` (String) The entity ref of the target of this relation.
- `type` (String) Type of the relation.

<a id="nestedatt--relations--target"></a>
### Nested Schema for `relations.target`

Read-Only:

- `kind` (String) The high level entity type being described.
- `name` (String) Name of the entity.
- `namespace` (String) Namespace that the target entity belongs to.



<a id="nestedatt--spec"></a>
### Nested Schema for `spec`

Read-Only:

- `children` (List of String) Children contains immediate child groups of this group in the hierarchy (whose parent field points to this group).
- `members` (List of String) Members contains the users that are members of this group.
- `parent` (String) Parent is the immediate parent group in the hierarchy, if any.
- `profile` (Attributes) Profile information about the group, mainly for display purposes. (see [below for nested schema](#nestedatt--spec--profile))
- `type` (String) The type of group.

<a id="nestedatt--spec--profile"></a>
### Nested Schema for `spec.profile`

Read-Only:

- `display_name` (String) A simple display name to present to users.
- `email` (String) Email where this entity can be reached.
- `picture` (String) A URL of an image that represents this entity.

