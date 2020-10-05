package bbva_fixtures

const CreateConsumerResponseSuccess = `{
	"user_id": "CO-64f149a1-c0ce-4980-9d18-6c71e88d9bd8",
	"address_id": "AD-cfc28966-c6cc-43fd-89e4-c109fe8e8b24",
	"addresses": [{
		"id": "AD-cfc28966-c6cc-43fd-89e4-c109fe8e8b24"
	}],
	"contact_id": ["CD-20d6f341-13bd-4ebd-8465-97aa3746baf7", "CD-98e7b987-0042-4e28-b232-b38dee34ef16"],
	"contacts": [{
		"id": "CD-20d6f341-13bd-4ebd-8465-97aa3746baf7"
	}, {
		"id": "CD-98e7b987-0042-4e28-b232-b38dee34ef16"
	}],
	"kyc": {
		"status": "APPROVED"
	},
	"kyc_notes": [{
		"code": "LV426"
	}, {
		"code": "NP100",
		"detail": "Detected the name matches the name record on the phone"
	}, {
		"code": "LV223"
	}],
	"digital_footprint": ["https://en.wikipedia.org/wiki/Application_programming_interface", "https://en.wikipedia.org/wiki/GIMP"]
}`

const CreateConsumerResponseInvalidAuthorizationToken = `{
	"errors": [{
		"code": "access_denied",
		"description": "Access denied"
	}]
}`
