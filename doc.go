// Package bazica converts a Gregorian birth time into a Ba-zi (Four Pillars of
// Destiny) chart: the year, month, day and hour pillars, the twelve life
// stages of each pillar, and the twelve luck pillars (Da Yun).
//
// The solar term and Lunar New Year tables are embedded in the module, so no
// data files are needed at runtime. GetBaziChart is safe for concurrent use.
package bazica
