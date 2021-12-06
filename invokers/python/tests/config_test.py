from pyfunc.config import Config
from pyfunc.constants import *

def test_load_default():
    cfg = Config()
    assert cfg.module_name == MODULE_NAME_DEFAULT
    assert cfg.function_name == FUNCTION_NAME_DEFAULT
    assert cfg.search_path == SEARCH_PATH_DEFAULT

def test_set_search_path():
    cfg = Config(search_path='./workspace')
    assert cfg.search_path == './workspace'