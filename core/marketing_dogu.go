package core

// MarketingDogu describes properties of a dogu which can be used to enhance the representation of a dogu in UI frontends.
//
// Example:
//
//	{
//		"ID": "0191b1e2-350f-7ecf-a0a8-6f2a2f0d3607",
//		"Deprecated": false,
//		"Namespace": "official",
//		"Name": "newdogu",
//		"DisplayName": "My New Dogu",
//		"Provider": [
//			{
//				"Slug": "cloudogu",
//				"Name": "Cloudogu GmbH",
//				"Logo": {
//					"ID": "34393ef9-a96d-4d1d-823e-0c17b10b762d",
//					"Title": "Cloudogu GmbH logo"
//				}
//			}
//		],
//		"Logo": {
//			"ID": "34393ef9-a96d-4d1d-823e-0c17b10b762d",
//			"Title": "Dogu icon"
//		},
//		"BackgroundImage": {
//			"ID": "18a5dc18-c202-40b2-9ee0-7f917ad82406",
//			"Title": "Screenshot of the main page of the new dogu"
//		},
//		"Descriptions": [
//			{
//				"Description": "Ein neues Dogu",
//				"LanguageCode": "de"
//			},
//			{
//				"Description": "A new dogu",
//				"LanguageCode": "en"
//			}
//		],
//		"ReleaseNotes": "https://example.com/release-notes"
//	}
type MarketingDogu struct {
	// ID contains the unique identifier of the dogu given by the CMS.
	//
	// Examples:
	//   - 0191b1e2-350f-7ecf-a0a8-6f2a2f0d3607
	//   - 34393ef9-a96d-4d1d-823e-0c17b10b762d
	//
	ID string
	// Deprecated indicates that this dogu is not recommended for new installations.
	// There will be no further development and updates.
	// In the Software Catalogue, deprecated dogus are marked with a warning sign.
	//
	Deprecated bool
	// Namespace contains the dogu's namespace without the name.
	//
	// Name together with Namespace delimited by a single forward slash "/" builds the dogu's full qualified name.
	//
	// Examples:
	//   - official
	//   - premium
	//
	// See also Dogu.Name for how dogu names are constructed.
	Namespace string
	// Name contains the dogu's simple name without the namespace.
	//
	// Name together with Namespace delimited by a single forward slash "/" builds the dogu's full qualified name.
	//
	// Examples:
	//   - redmine
	//   - confluence
	//   - newdogu
	//
	// See also Dogu.Name for how dogu names are constructed.
	Name string
	// Version defines the actual version of the dogu.
	//
	// The version follows the format from semantic versioning and additionally is split in two parts.
	// The application version and the dogu version.
	//
	// An example would be 1.7.8-1 or 2.2.0-4. The first part of the version (e.g. 1.7.8) represents the
	// version of the application (e.g. the nginx version in the nginx dogu). The second part represents the version
	// of the dogu and for an initial release it should start at 1 (e.g. 1.7.8-1).
	//
	// Example versions in the dogu.json:
	//  - 1.7.8-1
	//  - 2.2.0-4
	//
	Version string
	// PublishedAt is the date and time when the dogu was created.
	//
	// Examples:
	//   - 2024-10-16T07:49:34.738Z
	//   - 2019-05-03T13:31:48.612Z
	//
	PublishedAt string
	// DisplayName is the name of the dogu which is used in UI frontends to represent the dogu.
	//
	// Examples:
	//  - Jenkins CI
	//  - Backup & Restore
	//  - My New Dogu
	//
	DisplayName string
	// A provider is an entity that provides / maintains / develops the dogu.
	// The provider is a company in most cases.
	//
	// Example:
	// 	[{
	//		"Slug": "cloudogu",
	//		"Name": "Cloudogu GmbH",
	//		"Logo": {
	//			"ID": "34393ef9-a96d-4d1d-823e-0c17b10b762d",
	//			"Title": "Cloudogu GmbH logo"
	//		}
	//	}]
	//
	Provider []Provider
	// Logo contains information about the logo or icon of the dogu.
	//
	// Example:
	// 	{
	//		"ID": "34393ef9-a96d-4d1d-823e-0c17b10b762d",
	//		"Title": "Dogu Logo"
	//	}
	//
	Logo Image
	// BackgroundImage contains information about the background image of the dogu.
	//
	// Example:
	//	{
	//		"ID": "18a5dc18-c202-40b2-9ee0-7f917ad82406",
	//		"Title": "Screenshot of the main page of the new dogu"
	//	}
	//
	BackgroundImage Image
	// Descriptions contains a short explanation, what the dogu does in different languages.
	//
	// Example:
	// 	{
	//		"Description": "Ein neues Dogu",
	//		"LanguageCode": "de"
	//	},
	//	{
	//		"Description": "A new dogu",
	//		"LanguageCode": "en"
	//	}
	Descriptions []Translations
	// ReleaseNotes contains an URL to the release notes of the dogu.
	//
	// Example:
	//	- https://example.com/release-notes
	ReleaseNotes string
}

// Provider describes properties of the dogu's publisher.  These properties can be used to represent the provider in UI frontends so administrators and users can faster find appropriate support.
//
// Example:
//
//	{
//		"Slug": "cloudogu",
//		"Name": "Cloudogu GmbH",
//		"Logo": {
//			"ID": "34393ef9-a96d-4d1d-823e-0c17b10b762d",
//			"Title": "Cloudogu GmbH logo"
//		}
//	}
type Provider struct {
	// Slug contains the URL-ID of the provider which is part of an URL to access the providers page.
	//
	// Examples:
	//   - cloudogu
	//   - my-provider
	//
	Slug string
	// Name contains the name of the provider.
	//
	// Examples:
	//   - Cloudogu GmbH
	//   - Easy Software Ltd.
	//
	Name string
	// Logo contains information about the logo  of the provider.
	//
	// Example:
	// 	{
	//		"ID": "34393ef9-a96d-4d1d-823e-0c17b10b762d",
	//		"Title": "Provider Logo"
	//	}
	//
	Logo Image
}

// Image describes properties of an image which can be used to represent the image in UI frontends.
//
// Example:
//
//	{
//		"ID": "34393ef9-a96d-4d1d-823e-0c17b10b762d",
//		"Title": "Screenshot of the main page of the new dogu"
//	}
type Image struct {
	// ID contains the unique identifier of the image given by the CMS. The ID can be used to fetch the image from CMS.
	//
	// Examples:
	//   - 34393ef9-a96d-4d1d-823e-0c17b10b762d
	//	 - 18a5dc18-c202-40b2-9ee0-7f917ad82406
	//
	ID string
	// Title contains the title of the image which is used for accessibility reasons.
	//
	// Examples:
	//   - Screenshot of the main page of the new dogu
	//   - Logo of the provider
	//
	Title string
}

// Translations describes properties of a description of a dogu in a specific language. It can be used to represent the dogu's description in UI frontends
// in different languages.
//
// Example:
//
//	{
//		"Description": "MySQL - Relationale Datenbank",
//		"LanguageCode": "de"
//	},
//	{
//		"Description": "MySQL - Relational database",
//		"LanguageCode": "en"
//	}
type Translations struct {
	// Description contains a short explanation, what the dogu does in a specific language.
	//
	// Examples:
	//  - MySQL - Relationale Datenbank
	//  - Jenkins Continuous Integration Server
	//
	Description string
	// LanguageCode contains the ISO-639-1 code of the language in which the description is written.
	//
	// Examples:
	//   - de
	//   - en
	//
	LanguageCode string
}
