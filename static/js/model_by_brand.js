document.addEventListener("DOMContentLoaded", function () {
    const brandSelect = document.getElementById("brandSelect");
    const modelSelect = document.getElementById("modelSelect");

    if (!window.modelData) return;

    const allModels = window.modelData.allModels;
    const selectedModelID = window.modelData.selectedModelID;

    function populateModels(brandID) {
        if (!allModels || !Array.isArray(allModels)) return;
        modelSelect.innerHTML = '<option value="">-- Выберите модель --</option>';
        allModels.forEach(model => {
            if (model.MotorcycleBrandID === brandID) {
                const option = document.createElement("option");
                option.value = model.ID;
                option.textContent = model.model || model.MotoModel;
                if (model.ID === selectedModelID) {
                    option.selected = true;
                }
                modelSelect.appendChild(option);
            }
        });
    }

    const initialBrandID = parseInt(brandSelect.value);
    if (initialBrandID) {
        populateModels(initialBrandID);
    }

    brandSelect.addEventListener("change", function () {
        populateModels(parseInt(this.value));
    });
});
