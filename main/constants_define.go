package main

const UniqsTemplatePrefix      = "{{%"
const UniqsTemplateSuffix      = "%}}"
const UniqsTemplatePrefixValue = UniqsTemplatePrefix + "="
const UniqsTemplateLoopBegin   = UniqsTemplatePrefix + "loop.begin" + UniqsTemplateSuffix
const UniqsTemplateLoopEnd   = UniqsTemplatePrefix + "loop.end" + UniqsTemplateSuffix
const UniqsLoopLineIdx   = UniqsTemplatePrefix + "loop.lineIdx" + UniqsTemplateSuffix
