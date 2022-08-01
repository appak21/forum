let input, container, t;

TagsJoiner = '#get-tags'
allTags = document.querySelector(TagsJoiner)

input = document.querySelector('#hashtags');
container = document.querySelector('.tag-container');
let hashtagArray = [];

input.addEventListener('keyup', () => {
    if (event.which == 32 && input.value.trim().length > 0 ) {
      var text, p
      let newtags = input.value
      newtags.split(' ').forEach(element => {
        if (isValid(element)) {
          text = document.createTextNode(element);
          p = document.createElement('p');
          container.appendChild(p);
          p.appendChild(text);
          p.classList.add('tag');
          hashtagArray.push(element)
        }
      });

      input.value = '';

      allTags.value = hashtagArray.join(' ');
      
      let deleteTags = document.querySelectorAll('.tag');
      
      for(let i = 0; i < deleteTags.length; i++) {
        deleteTags[i].addEventListener('click', () => {
          container.removeChild(deleteTags[i]);
          let name = deleteTags[i].innerHTML
          hashtagArray = hashtagArray.filter(function(value, index, arr){
            return value != name
          });
          allTags.value = hashtagArray.join(' ');
        });
      };
    }
});

function isValid(tag) {
  if (tag=="" || hashtagArray.includes(tag)) {
    return false
  }
  if (tag.length > 30 || hashtagArray.length > 50) {
    return false
  }
  return true
}
