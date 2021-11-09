from flask import Request
from pyfunc.invoke import WrapFunction
import pytest

def test_wrap_function():
    is_called = False
    def testfunc(req):
        nonlocal is_called
        is_called = True
        return "called!"
    
    handler = WrapFunction(testfunc)
    handler()

    assert is_called

def test_unknown_name():
    is_called = False
    def testfunc(unknown):
        nonlocal is_called
        is_called = True
        return "called!"

    handler = WrapFunction(testfunc)
    with pytest.raises(AttributeError):
        handler()

@pytest.mark.skip(reason="Trying to run a method from Flask, causing it to fail")
def test_cloudevent_type():
    is_called = False
    def testfunc(event):
        nonlocal is_called
        is_called = True
        return "called!"

    handler = WrapFunction(testfunc)
    handler()

    assert is_called
