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

console.log(obj.test);

function named_function(){
    var a = 1;
    var b = 2;
    var c = a + b;
    var d = a + b;
    return c + d;
}