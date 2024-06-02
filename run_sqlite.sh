#!/bin/bash

BOLD='\033[1m'       #  ${BOLD} - жирный шрифт (интенсивный цвет)
BGMAGENTA='\033[45m'     #  ${BGMAGENTA}
LGREEN='\033[1;32m'     #  ${LGREEN}
NORMAL='\033[0m'      #  ${NORMAL} - все атрибуты по умолчанию

echo -e "${BOLD}${BGMAGENTA}${LGREEN} Running SQLite. ${NORMAL}"

cd internal/database || exit 1

sqlite3 ../storage_db/scheduler.db

cd ../..

echo -e "${BOLD}${BGMAGENTA}${LGREEN} SQLite work end. ${NORMAL}"

exit 0