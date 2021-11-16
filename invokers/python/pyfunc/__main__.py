from pyfunc import invoke

import sys

usage = f"""Usage:
  python -m framework
  python -m framework [search_path]
"""

search_path = "."
if len(sys.argv) == 2:
    search_path = sys.argv[1]
elif len(sys.argv) > 2:
    print(usage)
    exit(1)

invoke.main(search_path)
