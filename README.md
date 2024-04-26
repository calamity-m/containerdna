# containerdna

this is wip im learning leave me alone :)

Basic tool for viewing and validating things about a container's history

## Prerequisites

```
libbtrfs-dev pkg-config libgpgme-dev libdevmapper-dev
```

## Install

idk go get blabla

## Heritage

Complete a heritage check, which verifies that for every parent provided the child must originate from every single one.

This is done on a layer comparison basis. Given an example parent1 and parent2, the child must contain all layers of a
specified parent from its initial layer

	parent1 - -> layer0: A
	
	parent2 - -> layer0: A
		      -> layer1: AA

	child1  - -> layer0: A
		      -> layer1: AA
		      -> layer2: AAA

	child2  - -> layer0: A

With the default strict check:

	Child 1 is built from parent1 and parent2
	Child 2 is not built from parent1 and parent2, as it is lacking parent2's second layer.

With `--relaxed` flag supplied

	Child 1 and Child 2 are valid, as at least one parent is in their history

Running against remote registry:

	containerdna --relaxed --child docker://nginx --parent docker://nginx --parent docker://alpine --parent docker://ubuntu

Running against local daemon:

	containerdna --child docker-daemon:alpine:latest --parent docker-daemon:alpine:latest

