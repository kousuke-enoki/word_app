DICT_URL=https://github.com/scriptin/jmdict-simplified/releases/download/3.6.1%2B20250421122348/jmdict-eng-3.6.1+20250421122348.json.zip
ASSETS=backend/assets
$(ASSETS)/jmdict.json: ; @echo "Downloading JMdict..."
	@mkdir -p $(ASSETS)
	@curl -L -o /tmp/jmdict.zip $(DICT_URL)
	@unzip -d $(ASSETS) /tmp/jmdict.zip
	@mv $(ASSETS)/jmdict-eng-*.json $(ASSETS)/jmdict.json
	@rm /tmp/jmdict.zip
	@echo "JMdict ready."


download-dict: $(ASSETS)/jmdict.json
