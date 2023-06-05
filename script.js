document.addEventListener('DOMContentLoaded', function() {
    // Fetch the user list and display it on page load
    fetchUsers();

    // Add event listener to the create user form
    document.getElementById('create-user-form').addEventListener('submit', function(event) {
        event.preventDefault();

        // Get the form inputs
        var firstName = document.getElementById('first-name').value;
        var lastName = document.getElementById('last-name').value;
        var country = document.getElementById('country').value;

        // Create the user object
        var user = {
            first_name: firstName,
            last_name: lastName,
            country: country
        };

        // Send a POST request to create a new user
        createUser(user);
    });

    // Add event listener to the user list to handle delete and update actions
    document.getElementById('usersList').addEventListener('click', function(event) {
        var target = event.target;

        // Handle delete action
        if (target.classList.contains('delete-user')) {
            var userId = target.dataset.userId;
            deleteUser(userId);
        }

        // Handle update action
        if (target.classList.contains('update-user')) {
            var userId = target.dataset.userId;
            var firstName = prompt('Enter new first name:');
            var lastName = prompt('Enter new last name:');
            var country = prompt('Enter new country:');
            var user = {
                first_name: firstName,
                last_name: lastName,
                country: country
            };
            updateUser(userId, user);
        }
    });
});

function fetchUsers() {
    // Send a GET request to retrieve the list of users
    fetch('http://localhost:8080/users')
        .then(response => response.json())
        .then(data => {
            // Clear the user list
            var userList = document.getElementById('usersList');
            userList.innerHTML = '';

            // Iterate over the users and add them to the list
            data.forEach(user => {
                var listItem = document.createElement('li');
                listItem.textContent = user.first_name + ' ' + user.last_name;

                var deleteButton = document.createElement('button');
                deleteButton.textContent = 'Delete';
                deleteButton.classList.add('delete-user');
                deleteButton.dataset.userId = user.id;

                var updateButton = document.createElement('button');
                updateButton.textContent = 'Update';
                updateButton.classList.add('update-user');
                updateButton.dataset.userId = user.id;

                listItem.appendChild(deleteButton);
                listItem.appendChild(updateButton);

                userList.appendChild(listItem);
            });
        })
        .catch(error => {
            console.error('Error:', error);
        });
}

function createUser(user) {
    // Send a POST request to create a new user
    fetch('http://localhost:8080/users', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(user)
    })
    .then(response => response.json())
    .then(data => {
        // Clear the form inputs
        document.getElementById('first-name').value = '';
        document.getElementById('last-name').value = '';
        document.getElementById('country').value = '';

        // Fetch the updated user list and display it
        fetchUsers();
    })
    .catch(error => {
        console.error('Error:', error);
    });
}

function deleteUser(userId) {
    // Send a DELETE request to delete the user
    fetch(`http://localhost:8080/users/${userId}`, {
        method: 'DELETE'
    })
    .then(response => response.json())
    .then(data => {
        // Fetch the updated user list and display it
        fetchUsers();
    })
    .catch(error => {
        console.error('Error:', error);
    });
}

function updateUser(userId, user) {
    // Send a PUT request to update the user
    fetch(`http://localhost:8080/users/${userId}`, {
        method: 'PUT',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(user)
    })
    .then(response => response.json())
    .then(data => {
        // Fetch the updated user list and display it
        fetchUsers();
    })
    .catch(error => {
        console.error('Error:', error);
    });
}
