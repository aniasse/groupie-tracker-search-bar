var selectOptions = document.querySelectorAll('.select-option');
selectOptions.forEach(function(select) {
    select.onchange = function() {
        var selectedOption = select.options[select.selectedIndex];
        var url = selectedOption.value;

        window.location.href=url;
    };
});