from django.shortcuts import render
from django.http import HttpResponse,HttpResponseNotFound
from .forms import AddPostForm
from .models import User
#работает через жепу тк по урлу просит ф-цию представления(одну)
def homepage(request):
   
    return render(request,"lendpage/index.html")

#отображает всратые поля(можно поправить), ебаное поле email
def display_form(request):
    if request.method == "POST":
        form = AddPostForm(request.POST)
        if form.is_valid():
            try:
                User.objects.create(**form.cleaned_data)
                #return redirect('homepage')
            except:
                form.add_error(None, "Error message")
    else:
        form=AddPostForm()

    
    data = {
        'form':form,
    }
    return render(request,"lendpage/index.html", data)

#отключив режим debug сосу бербу(не дружественный вид) ф-ция не робит(прописал в урлах)
def page_not_found(request,exeption):
    return HttpResponseNotFound("Страница не найдена")