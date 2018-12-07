var b = document.getElementById('a').value;
eval(b);

(function () {
    var x = 9;
    var a = x * b;
})();

var obj = {
    prop: "test"
};

if(typeof obj == 'object'){
    alert(obj.prop)
}

console.log(obj.test)