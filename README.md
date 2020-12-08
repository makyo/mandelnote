---
title: Mandelnote
author: Madison Scott-Clary
description: A fractal outlining tool for the terminal
---

# Mandelnote

A fractal outlining tool for the terminal.

## Usage

    mandelnote <document>

Once in there, `ctrl+H` will get you the help. This document is a valid Mandelnote file.

# About

## Goals

The goal of Mandelnote is to provide a unique outlining and writing experience based loosely on the snowflake method: acts contain chapters contain scenes, which can then be written and then merged and promoted until one has a completed work. Of course, it's more general than that, it doesn't rely on the idea of acts/chapters/scenes, of course.

Additionally, the native file format is just markdown with a YAML block, suitable for using in Jekyll, Hugo, and the like.

## Rationale

Madison Scott-Clary found myself writing pretty often in a variation of what's called the snowflake method, where you start with an idea, then come up with an outline of acts, then come up with an outline of chapters, then come up with an outline of scenes. Then you write the scenes and flatten them into chapters, then flatten the chapters into acts, then flatten the acts into the final product.

To that end, she decided to poke at making an editor for just that. It stores each of those elements in cards in a notebook which one can move between set titles, set contents, etc. The format that it uses is just Markdown: cards are a header of any depth followed by the text of the body.

### For example...

Although it was written before Mandelnote's inception, the file format was used with Madison Scott-Clary's book [Qoheleth](https://qoheleth.makyo.ink).

## Contributors

* [Madison Scott-Clary](https://makyo.is)

### Contributing

Always open to help with the code!
