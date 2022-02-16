package functions;

import java.util.function.Function;

import functions.models.Employee;
import functions.models.Person;

/*
This class demonstrates the definition of a function called "hire".
This function can be accessed by targetting the "/hire" path while
providing the correct data:
    {
        "specversion" : "1.0",
        "type" : "org.springframework",
        "source" : "https://spring.io/",
        "id" : "A234-1234-1234",
        "datacontenttype" : "application/json",
        "data": {
            "firstName": "John",
            "lastName": "Doe"
        }
    }
If this is the only function defined, it may be accessed via "/"
path.
*/
public class Hire implements Function<Person, Employee> {
    @Override
    public Employee apply(Person person) {
        System.out.printf("Person: first(%s) last(%s)\n", person.getFirstName(), person.getLastName());
        Employee employee = new Employee(person);
        System.out.printf("Employee: %s\n", employee.getMessage());
        return employee;
    }
}
