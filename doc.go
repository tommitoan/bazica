// Package bazica converts a Gregorian birth time into a Ba-zi (Four Pillars of
// Destiny) chart: the year, month, day and hour pillars, the twelve life
// stages of each pillar, and the twelve luck pillars (Da Yun).
//
// # Analysis
//
// Every chart also carries derived facts under an "analysis" key, next to the
// fields that existed in v1.3.0, which are unchanged:
//
//   - each natal pillar: its Ten God against the Day Master, the elements,
//     Nayin, life stages, hidden stems with their Ten Gods and stages, whether
//     its branch is void, a heaven-clash/earth-clash flag and its stars;
//   - each luck pillar: nominal age range, Ten God, Nayin, life stage, hidden
//     stems and the clash flag;
//   - the chart: the Day Master, the void branches, the conception (Thai Nguyen),
//     gestation (Thai Tuc) and life palace pillars, and unweighted element counts.
//
// GetAnnualPillars returns the yearly pillars with the same terms and flags.
//
// Every named concept is a model.LocalizedTerm: a stable Code, an English label
// and a Vietnamese label. Switch on Code, not on a label. Fields that are
// defined but not computed (strength and Useful God) are present and null.
//
// # Conventions
//
//   - Ages are nominal: the year of the year pillar counts as age 1. The year
//     pillar's year is the Lunar New Year year, which for a January birth
//     before the Lunar New Year is the year before the civil birth year.
//   - Stars are listed only when their rules have been verified against
//     reference charts; other stars are never reported.
//   - The life palace puts the Rat and Ox branches before the Tiger when it
//     derives the stem, which differs from the classical month order.
//   - The heaven-clash/earth-clash flag is directional: a cell is flagged when
//     its stem overcomes another natal pillar's stem with the same polarity and
//     their branches are opposite.
//
// The solar term and Lunar New Year tables are embedded in the module, so no
// data files are needed at runtime. GetBaziChart and GetAnnualPillars are safe
// for concurrent use.
package bazica
