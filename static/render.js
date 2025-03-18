document.addEventListener("DOMContentLoaded", async () => {
    const listContainer = document.getElementById("listContainer");
    const taskContainer = document.getElementById("taskContainer");
    const listTitle = document.getElementById("listTitle");
    const addTaskBtn = document.getElementById("addTaskBtn");
    const deleteListBtn = document.getElementById("deleteListBtn");

    let lists = [];
    let selectedList = null;

    // 🔹 Récupérer les listes depuis l'API au démarrage
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

    // 🔹 Récupérer les tâches d'une liste
    async function fetchTasks(listId) {
        try {
            const response = await fetch(`/list/${listId}/tasks`);
            selectedList.tasks = await response.json();
            renderTasks();
        } catch (error) {
            console.error("Erreur lors du chargement des tâches :", error);
        }
    }

    // 🔹 Afficher la liste des listes
    function renderLists() {
        listContainer.innerHTML = "";
        lists.forEach((list, index) => {
            const li = document.createElement("li");
            li.textContent = list.name;
            li.addEventListener("click", () => selectList(index));
            listContainer.appendChild(li);
        });
    }

    // 🔹 Sélectionner une liste et charger ses tâches
    async function selectList(index) {
        selectedList = lists[index];
        listTitle.textContent = `📌 ${selectedList.name}`;
        addTaskBtn.classList.remove("hidden");
        deleteListBtn.classList.remove("hidden");
        await fetchTasks(selectedList.id);
    }

    // 🔹 Ajouter une nouvelle liste
    async function addList() {
        const listName = prompt("Nom de la nouvelle liste :");
        if (!listName) return;
        try {
            const response = await fetch("/api/lists", {
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

    // 🔹 Ajouter une nouvelle tâche
    async function addTask() {
        const taskName = prompt("Nouvelle tâche :");
        if (!taskName || !selectedList) return;
        try {
            const response = await fetch(`/list/${selectedList.id}/task`, {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ description: taskName }),
            });
            const newTask = await response.json();
            selectedList.tasks.push(newTask);
            renderTasks();
        } catch (error) {
            console.error("Erreur lors de l'ajout de la tâche :", error);
        }
    }

    // 🔹 Cocher/Décocher une tâche (avec mise à jour en BDD)
    async function toggleTask(taskId, checked) {
        try {
            await fetch(`/task/${taskId}/check`, { method: "PATCH" });
            selectedList.tasks = selectedList.tasks.map(task => 
                task.id === taskId ? { ...task, checked: !checked } : task
            );
            renderTasks();
        } catch (error) {
            console.error("Erreur lors de la mise à jour de la tâche :", error);
        }
    }

    // 🔹 Supprimer une liste
    async function deleteList() {
        if (!selectedList || !confirm("Supprimer cette liste ?")) return;
        try {
            await fetch(`/api/list/${selectedList.id}`, { method: "DELETE" });
            lists = lists.filter(list => list.id !== selectedList.id);
            selectedList = lists.length > 0 ? lists[0] : null;
            renderLists();
            renderTasks();
        } catch (error) {
            console.error("Erreur lors de la suppression :", error);
        }
    }

    // 🔹 Drag & Drop pour réorganiser les tâches (avec envoi en BDD)
    function setupDragAndDrop() {
        let draggedItem = null;

        taskContainer.addEventListener("dragstart", (e) => {
            draggedItem = e.target;
            e.dataTransfer.setData("text/plain", draggedItem.dataset.id);
        });

        taskContainer.addEventListener("dragover", (e) => {
            e.preventDefault();
        });

        taskContainer.addEventListener("drop", async (e) => {
            e.preventDefault();
            const droppedOn = e.target.closest(".task-item");
            if (!draggedItem || !droppedOn || draggedItem === droppedOn) return;

            const draggedTaskId = draggedItem.dataset.id;
            const droppedOnTaskId = droppedOn.dataset.id;

            // Échanger les tâches dans la BDD
            try {
                await fetch(`/task/${draggedTaskId}/swapOrder`, {
                    method: "PATCH",
                    headers: { "Content-Type": "application/json" },
                    body: JSON.stringify({ taskToSwapId: droppedOnTaskId }),
                });
                await fetchTasks(selectedList.id); // Recharger les tâches après l'échange
            } catch (error) {
                console.error("Erreur lors du réarrangement des tâches :", error);
            }
        });
    }

    document.getElementById("addListBtn").addEventListener("click", addList);
    addTaskBtn.addEventListener("click", addTask);
    deleteListBtn.addEventListener("click", deleteList);

    await fetchLists(); // Charge les listes au démarrage
    setupDragAndDrop(); // Active le drag & drop
});