package regolint.rules.imports.banned

metadata := {
	"id": "IMP001",
	"severity": "error",
	"description": "Prevents use of unsafe package",
}

deny contains violation if {
	some imp in input.imports
	imp.path == "unsafe"

	violation := {
		"message": "Import of banned package 'unsafe'",
		"position": imp.position,
		"rule": metadata.id,
	}
}
