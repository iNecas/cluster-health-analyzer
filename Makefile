MERMAID_PATH="mmdc"

MERMAID_OUT_FILES=$(patsubst %.mmd, %.svg, $(shell find docs -name "*.mmd"))


.PHONY: mermaid
mermaid: $(MERMAID_OUT_FILES)

%.svg: %.mmd
	$(MERMAID_PATH) -i $< -o $@
