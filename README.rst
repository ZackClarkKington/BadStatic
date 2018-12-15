===============================================================================
BadStatic
===============================================================================

*An attempt at writing a static analyzer for JavaScript*

**Specification of rules**:

An example showing how to specify a set of rules for the analyzer to follow can be found in the file test.json.

Rules have different types, the types currently implemented are:

- **Expression**, for use in monitoring and addressing usage of a specific expression (such as eval, as seen in test.js)
- **PropertyDoesNotExist**, for finding attempts to access undefined properties of an object

Each rule should have an associated action.

**Specification of actions**:

Actions have two properties:

- **info** - Contains an information message that will be printed when the action is executed
- **type** - Can have two values, *fail* or *warn*, *fail* will exit immediately and no further analysis will take place whereas *warn* will just print the information message defined in the action
