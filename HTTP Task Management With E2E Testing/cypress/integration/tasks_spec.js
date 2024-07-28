describe('Task Management', () => {
    let token;

    before(() => {
        // Register a new user
        cy.request('POST', '/register', {
            username: 'testuser',
            password: 'password'
        });

        // Log in to get the token
        cy.request('POST', '/login', {
            username: 'testuser',
            password: 'password'
        }).then((response) => {
            token = response.body.token;
        });
    });

    it('should create a new task', () => {
        cy.request({
            method: 'POST',
            url: '/tasks',
            headers: {
                Authorization: token
            },
            body: {
                description: 'New Task'
            }
        }).then((response) => {
            expect(response.status).to.eq(201);
            expect(response.body.description).to.eq('New Task');
        });
    });

    it('should list all tasks', () => {
        cy.request({
            method: 'GET',
            url: '/tasks',
            headers: {
                Authorization: token
            }
        }).then((response) => {
            expect(response.status).to.eq(200);
            expect(response.body).to.have.length(1);
        });
    });

    it('should update a task', () => {
        cy.request({
            method: 'GET',
            url: '/tasks',
            headers: {
                Authorization: token
            }
        }).then((response) => {
            const taskId = response.body[0].id;
            cy.request({
                method: 'PUT',
                url: `/tasks/${taskId}`,
                headers: {
                    Authorization: token
                },
                body: {
                    description: 'Updated Task',
                    completed: true
                }
            }).then((response) => {
                expect(response.status).to.eq(200);
                expect(response.body.description).to.eq('Updated Task');
                expect(response.body.completed).to.eq(true);
            });
        });
    });

    it('should delete a task', () => {
        cy.request({
            method: 'GET',
            url: '/tasks',
            headers: {
                Authorization: token
            }
        }).then((response) => {
            const taskId = response.body[0].id;
            cy.request({
                method: 'DELETE',
                url: `/tasks/${taskId}`,
                headers: {
                    Authorization: token
                }
            }).then((response) => {
                expect(response.status).to.eq(204);
            });
        });
    });
});
