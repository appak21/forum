let input, hashtagArray, container, t;

TagsJoiner = '#get-tags'
allTags = document.querySelector(TagsJoiner)

input = document.querySelector('#hashtags');
container = document.querySelector('.tag-container');
hashtagArray = [];

input.addEventListener('keyup', () => {
  hashtagArray = [];
    if (event.which == 32 && input.value.length > 0) {
      var text = document.createTextNode(input.value);
      var p = document.createElement('p');
      container.appendChild(p);
      p.appendChild(text);
      p.classList.add('tag');
      hashtagArray.push(input.value)
      input.value = '';
      
      let deleteTags = document.querySelectorAll('.tag');
      
      for(let i = 0; i < deleteTags.length; i++) {
        deleteTags[i].addEventListener('click', () => {
          container.removeChild(deleteTags[i]);
        });
      }
      
      let newTagArr = []
      let remtags = deleteTags;
      for(let i = 0; i < remtags.length; i++) {
        newTagArr.push(remtags[i].innerHTML)
      }
      allTags.value = newTagArr.join(' ')
    }
});

