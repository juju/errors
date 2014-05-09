/*
package arrar provides an easy way to annotate errors without losing the
orginal error context.

The package does not export any error types that are expected to be
composed, but instead operates primarily by annotating other error types.

The exported Errorf function is designed to replace the fmt.Errorf
function. The same underlying error is there, but the package also records
the location at which the error was created.

A primary use case for this library is to add extra context any time an
error is returned from a function.

    if err := SomeFunc(); err != nil {
	    return err
	}

This instead becomes:

    if err := SomeFunc(); err != nil {
	    return arrar.Trace(err)
	}

which just records the file, line and function, or

    if err := SomeFunc(); err != nil {
	    return arrar.Annotate(err, "more context")
	}

which adds annotation to the error.

You are able to get the underlying error back from the wrapped error by calling

	lastError := arrar.LastError(err)

Often when you want to check to see if an error is of a particular type, a
helper function is exported by the package that returned the error, like the
`os` package.  Since arrar wraps the error, you cannot directly test the
error, but instead need to pass the checking function through.

	arrar.Check(err, os.IsNotExist)

The result of the Error() call on the annotated error is the annotations
joined with commas, followed by a colon, then the result of the Error() method
for the underlying error.

	err := arrar.Errorf("original")
	err = arrar.Annotatef("context")
	err = arrar.Annotatef("more context")
	err.Error() -> "more context, context: original"

Obviously recording the file, line and functions is not very useful if you
cannot get them back out again.

	arrar.DefaultErrorStack(err)

will return something like:

	four [four@github.com/juju/arrar/test_functions_test.go:32]
	translated [transthree@github.com/juju/arrar/test_functions_test.go:28]
	two: one [two@github.com/juju/arrar/test_functions_test.go:16]

where the format of the line is: annotation: error [func@file:line]. The most
recently annotated or wrapped message is shown at the top, and the first
annotating call last.

If you are creating the errors, you can simply call:

	arrar.Errorf("format just like fmt.Errorf")

This function will return an error that contains the annotation stack and
records the file, line and function from the place where the error is created.

Sometimes when responding to an error you want to return a more specific error
for the situation.

    if err := FindField(field); err != nil {
	    return arrar.Wrap(err, NotFoundError(field))
	}

This returns an error where the complete error stack is still available, and
arrar.LastError will return the NotFoundError.

You are able to get a slice all of the actual error values captured using

	arrar.ErrorStack()

The original error is the first value, and the error from the most recent Wrap
call is last.

CAVEAT:

gccgo currently returns mangled names for method calls which are not easy to
demangle, so function names are not recorded if using gccgo.

*/
package errors
