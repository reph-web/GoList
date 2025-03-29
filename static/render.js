document.addEventListener("DOMContentLoaded", async () => {
    const listContainer = document.getElementById("listContainer");
    const taskContainer = document.getElementById("taskContainer");
    const listTitle = document.getElementById("listTitle");
    const addTaskBtn = document.getElementById("addTaskBtn");
    const deleteListBtn = document.getElementById("deleteListBtn");

    let lists = [];
    let selectedList = null;

    async function fetchLists() {
        try {
            const response = await fetch("/api/lists");
            lists = await response.json();
            if (lists.length > 0) {
                selectList(0);
            } else {
                listTitle.textContent = "📝 Aucune liste disponible";
            }
            renderLists();
        } catch (error) {
            console.error("Erreur lors du chargement des listes :", error);
        }
    }

    async function fetchTasks(listId) {
        try {
            const response = await fetch(`/api/list/${listId}/tasks`);
            selectedList.Tasks = await response.json();
            renderTasks(selectedList.Tasks);
        } catch (error) {
            console.error("Erreur lors du chargement des tâches :", error);
        }
    }
    function renderLists() {
        listContainer.innerHTML = "";
        lists.forEach((list, index) => {
            const li = document.createElement("li");
            li.textContent = list.Name;
            li.addEventListener("click", () => selectList(index));
            listContainer.appendChild(li);
        });
    }

    
    function renderTasks(selectedListTasks) {
        taskContainer.innerHTML = "";
        selectedListTasks.forEach((task) => {
            const li = document.createElement("li");
            const checkbox = document.createElement("input");
            checkbox.type = "checkbox";
            checkbox.checked = task.Checked;
            checkbox.addEventListener("change", () => toggleTask(task.ID));
            li.appendChild(checkbox);

            const description = document.createElement("span");
            description.textContent = task.Description;
            li.appendChild(description);

            const deleteButton = document.createElement("button");
            deleteButton.textContent = "❌";
            deleteButton.addEventListener("click", () => deleteTask(task.ID));
            li.appendChild(deleteButton);

            taskContainer.appendChild(li);
        });
    }

    async function selectList(index) {
        selectedList = lists[index];
        listTitle.textContent = `📌 ${selectedList.Name}`;
        addTaskBtn.classList.remove("hidden");
        deleteListBtn.classList.remove("hidden");
        await fetchTasks(selectedList.ID);
    }

    async function addList() {
        const listName = prompt("Nom de la nouvelle liste :");
        if (!listName) return;
        try {
            const response = await fetch("/api/list", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ name: listName }),
            });
            const newList = await response.json();
            lists.push(newList);
            renderLists();
        } catch (error) {
            console.error("Erreur lors de la création de la liste :", error);
        }
    }

    async function addTask() {
        const taskName = prompt("Nouvelle tâche :");
        if (!taskName || !selectedList) return;
        try {
            const response = await fetch(`/api/list/${selectedList.ID}/task`, {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ description: taskName }),
            });
            const newTask = await response.json();
            selectedList.Tasks.push(newTask);
            renderTasks();
        } catch (error) {
            console.error("Erreur lors de l'ajout de la tâche :", error);
        }
    }

    async function toggleTask(taskId) {
        try {
            await fetch(`api/task/${taskId}/check`, { method: "PATCH" });
        } catch (error) {
            console.error("Erreur lors de la mise à jour de la tâche :", error);
        }
    }

    async function deleteList() {
        if (!selectedList || !confirm("Supprimer cette liste ?")) return;
        try {
            await fetch(`/api/list/${selectedList.ID}`, { method: "DELETE" });
            lists = lists.filter(list => list.ID !== selectedList.ID);
            selectedList = lists.length > 0 ? lists[0] : null;
            renderLists();
            renderTasks();
        } catch (error) {
            console.error("Erreur lors de la suppression :", error);
        }
    }

    document.getElementById("addListBtn").addEventListener("click", addList);
    addTaskBtn.addEventListener("click", addTask);
    deleteListBtn.addEventListener("click", deleteList);

    await fetchLists();
});