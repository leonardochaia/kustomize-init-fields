# Kustomize Field Init

When using Kustomize, there is no way to "append an item to a sequence,
creating the sequence if it doesnt exist".

This plugin initializes common container fields so that I don't hit that issue
while patching those fields.

Unfortunately, I can't use Strategic Merge Patch since I don't know the names of the
containers beforehand.