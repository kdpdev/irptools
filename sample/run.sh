#!/bin/bash

export PATH=$PATH:.
irptools.exe -cmd=parse -cfg=cfg_parse.json
irptools.exe -cmd=filter -cfg=cfg_filter.json
irptools.exe -cmd=export_fz -cfg=cfg_export_fz.json
