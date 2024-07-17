from django.shortcuts import render
from django.http import HttpResponse,HttpResponseNotFound
from .forms import AddPostForm
from .models import User

def homepage(request):
   
    return render(request,"lendpage/index.html")


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


def page_not_found(request,exeption):
    return HttpResponseNotFound("Страница не найдена")