package components

// PhysicsComponent contains physics data for an entity
type PhysicsComponent struct {
	ForceDown  float64
	ForceRight float64
	ForceUp    float64
	ForceLeft  float64
	Blocked    bool
}
